package elastic

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/db/srv/engine"
	data "github.com/kazoup/platform/db/srv/engine/elastic/data"
	db "github.com/kazoup/platform/db/srv/proto/db"
	search_proto "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/kazoup/platform/structs/globals"
	lib "github.com/mattbaird/elastigo/lib"
	"golang.org/x/net/context"
)

type elastic struct {
	conn                 *lib.Conn
	bulk                 *lib.BulkIndexer
	filesChannel         chan *crawler.FileMessage
	slackUsersChannel    chan *crawler.SlackUserMessage
	slackChannelsChannel chan *crawler.SlackChannelMessage
	crawlerFinished      chan *crawler.CrawlerFinishedMessage
	esMapping            *[]byte // For files
	esSettings           *[]byte // For files index
}

func init() {
	es_mapping, err := data.Asset("data/es_mapping_files_new.json")
	if err != nil {
		// Asset was not found.
		log.Fatal(err)
	}
	es_settings, err := data.Asset("data/es_settings.json")
	if err != nil {
		// Asset was not found.
		log.Fatal(err)
	}

	engine.Register(&elastic{
		filesChannel:         make(chan *crawler.FileMessage),
		slackUsersChannel:    make(chan *crawler.SlackUserMessage),
		slackChannelsChannel: make(chan *crawler.SlackChannelMessage),
		crawlerFinished:      make(chan *crawler.CrawlerFinishedMessage),
		esMapping:            &es_mapping,
		esSettings:           &es_settings,
	})
}

// Init elastic db
func (e *elastic) Init() error {
	e.conn = lib.NewConn()
	e.conn.SetHosts([]string{"localhost:9200"}) //TODO: replace for enterprise version, get flag
	e.bulk = e.conn.NewBulkIndexerErrors(100, 5)
	e.bulk.BulkMaxDocs = 100000
	e.bulk.Start()

	if err := indexer(e); err != nil {
		return err
	}

	if err := enricher(e); err != nil {
		return err
	}

	return nil
}

// Create record
func (e *elastic) Create(req *db.CreateRequest) (*db.CreateResponse, error) {
	_, err := e.conn.Index(req.Index, req.Type, req.Id, nil, req.Data)

	return &db.CreateResponse{}, err
}

// Subscribe to crawler file messages
func (e *elastic) SubscribeFiles(ctx context.Context, msg *crawler.FileMessage) error {
	e.filesChannel <- msg

	return nil
}

// Subscribe to crawler file messages
func (e *elastic) SubscribeSlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error {
	e.slackUsersChannel <- msg

	return nil
}

// Subscribe to crawler file messages
func (e *elastic) SubscribeSlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error {
	e.slackChannelsChannel <- msg

	return nil
}

// Subscribe to crawler finished message
func (e *elastic) SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	e.crawlerFinished <- msg

	return nil
}

// Read record from ES
// Since we can't get the record bye id from aliases we need to use search request
// This should return single ID as all files should have unique ID's as we seting them up based on unique path MD5
// Method will worl on any index and alias as long ID's are unique
// ES query
//{
//  "query": {
//    "term": {
//      "id":"7041aef7f8e9fbd5a0379d04e21d3ecf"
//    }
//  }
//}
type SearchRequestByID struct {
	Query *Query `json:"query"`
}
type Query struct {
	Term *Term `json:"term"`
}
type Term struct {
	ID string `json:"id"`
}

func (e *elastic) Read(req *db.ReadRequest) (*db.ReadResponse, error) {

	// build query
	t := &Term{
		ID: req.Id,
	}
	q := &Query{
		Term: t,
	}
	sr := &SearchRequestByID{
		Query: q,
	}
	//Marshal query
	bq, err := json.Marshal(sr)
	if err != nil {
		return &db.ReadResponse{}, err
	}
	//Query string
	query := string(bq)
	//Search request
	out, err := e.conn.Search(req.Index, req.Type, nil, query)
	if err != nil {
		return &db.ReadResponse{}, err
	}
	v := out.Hits.Hits
	// hmmm hacky FIXME
	if out.Hits.Total == 0 {
		return &db.ReadResponse{}, errors.New("Got 0 results")

	}
	if out.Hits.Total != 1 {
		return &db.ReadResponse{}, errors.New("Got more then one results")

	}
	// Now we should have only one result
	data, err := v[0].Source.MarshalJSON()

	if err != nil {
		return &db.ReadResponse{}, err
	}
	s, err := engine.TypeFactory(req.Type, string(data))
	if err != nil {
		return &db.ReadResponse{}, err
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return &db.ReadResponse{}, err
	}

	response := &db.ReadResponse{
		Result: string(data),
	}
	return response, nil
}

// Update record
func (e *elastic) Update(req *db.UpdateRequest) (*db.UpdateResponse, error) {
	_, err := e.conn.Index(req.Index, req.Type, req.Id, nil, req.Data)

	return &db.UpdateResponse{}, err
}

// Delete record
func (e *elastic) Delete(req *db.DeleteRequest) (*db.DeleteResponse, error) {
	_, err := e.conn.Delete(req.Index, req.Type, req.Id, nil)

	return &db.DeleteResponse{}, err
}

// Search ES index
func (e *elastic) Search(req *db.SearchRequest) (*db.SearchResponse, error) {
	var results []interface{}

	eQuery := ElasticQuery{
		Index:    req.Index,
		Term:     req.Term,
		From:     req.From,
		Size:     req.Size,
		Category: req.Category,
		Url:      req.Url,
		Depth:    req.Depth,
		Type:     req.FileType,
	}
	query, err := eQuery.Query()

	if err != nil {
		return &db.SearchResponse{}, err
	}

	out, err := e.conn.Search(req.Index, req.Type, nil, query)
	if err != nil {
		return &db.SearchResponse{}, err
	}

	for _, v := range out.Hits.Hits {
		data, err := v.Source.MarshalJSON()
		if err != nil {
			return &db.SearchResponse{}, err
		}
		s, err := engine.TypeFactory(req.Type, string(data))
		if err != nil {

			return &db.SearchResponse{}, err
		}
		if err := json.Unmarshal(data, &s); err != nil {
			return &db.SearchResponse{}, err
		}
		results = append(results, s)
	}

	info := gabs.New()
	info.Set(out.Hits.Total, "total")

	b, err := json.Marshal(results)
	if err != nil {
		return &db.SearchResponse{}, err
	}

	return &db.SearchResponse{
		Result: string(b),
		Info:   info.String(),
	}, nil
}

// CreateIndexWithSettings creates an ES index with settings
func (e *elastic) CreateIndexWithSettings(req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error) {
	var settingsMap map[string]interface{}

	exists, err := e.conn.IndicesExists(req.Index)
	if err != nil {
		return &db.CreateIndexWithSettingsResponse{}, err
	}

	if !exists {
		if err := json.Unmarshal(*e.esSettings, &settingsMap); err != nil {
			return &db.CreateIndexWithSettingsResponse{}, err
		}

		_, err := e.conn.CreateIndexWithSettings(req.Index, settingsMap)
		if err != nil {
			return &db.CreateIndexWithSettingsResponse{}, err
		}
	}

	return &db.CreateIndexWithSettingsResponse{}, nil
}

// PutMappingFromJSON puts a mapping into ES
func (e *elastic) PutMappingFromJSON(req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error) {
	var clusterHealth lib.ClusterHealthResponse
	var err error

	// Check for cluster health, continue when changes from red to safer (yellow / green)
	// http://xbib.org/elasticsearch/2.1.1/apidocs/org/elasticsearch/indices/IndexPrimaryShardNotAllocatedException.html
	clusterHealth, err = e.conn.Health(req.Index)
	if err != nil {
		return &db.PutMappingFromJSONResponse{}, err
	}

	for clusterHealth.Status == "red" {
		clusterHealth, err = e.conn.Health(req.Index)
		if err != nil {
			return &db.PutMappingFromJSONResponse{}, err
		}
	}

	if len(req.Type) == 0 {
		return nil, errors.New("document type required")
	}

	if _, err := e.conn.CloseIndex(req.Index); err != nil {
		return &db.PutMappingFromJSONResponse{}, err
	}

	if err := e.conn.PutMappingFromJSON(req.Index, req.Type, *e.esMapping); err != nil {
		return nil, err
	}

	if _, err := e.conn.OpenIndex(req.Index); err != nil {
		return &db.PutMappingFromJSONResponse{}, err
	}

	return &db.PutMappingFromJSONResponse{}, nil
}

// Status elasticsearch cluster
func (e *elastic) Status(req *db.StatusRequest) (*db.StatusResponse, error) {
	clusterState, err := e.conn.ClusterState(lib.ClusterStateFilter{
		FilterNodes:        true,
		FilterRoutingTable: true,
		FilterMetadata:     true,
		FilterBlocks:       true,
	})

	b, err := json.Marshal(clusterState)
	if err != nil {
		return &db.StatusResponse{}, err
	}

	response := &db.StatusResponse{
		Status: string(b),
	}

	return response, err
}

// AddAlias to assign indexes (aliases) per datasource
func (e *elastic) AddAlias(req *db.AddAliasRequest) (*db.AddAliasResponse, error) {
	_, err := e.conn.AddAlias(req.Index, req.Alias)
	if err != nil {
		return nil, err
	}

	return &db.AddAliasResponse{}, nil
}

// DeleteIndex from ES
func (e *elastic) DeleteIndex(req *db.DeleteIndexRequest) (*db.DeleteIndexResponse, error) {
	_, err := e.conn.DeleteIndex(req.Index)
	if err != nil {
		return nil, err
	}

	return &db.DeleteIndexResponse{}, nil
}

// DeleteAlias from ES
func (e *elastic) DeleteAlias(req *db.DeleteAliasRequest) (*db.DeleteAliasResponse, error) {
	_, err := e.RemoveAlias(req.Index, req.Alias)
	if err != nil {
		return nil, err
	}

	return &db.DeleteAliasResponse{}, nil
}

// RenameAlias from ES
func (e *elastic) RenameAlias(req *db.RenameAliasRequest) (*db.RenameAliasResponse, error) {
	var err error

	_, err = e.RemoveAlias(req.Index, req.OldAlias)
	if err != nil {
		return nil, err
	}

	_, err = e.AddAlias(&db.AddAliasRequest{
		Index: req.Index,
		Alias: req.NewAlias,
	})
	if err != nil {
		return nil, err
	}

	return &db.RenameAliasResponse{}, nil
}

// Aggregate allow us to query for aggs in ES
func (e *elastic) Aggregate(req *search_proto.AggregateRequest) (*search_proto.AggregateResponse, error) {
	eQuery := ElasticQuery{
		Term:     req.Filters.Term,
		Category: req.Filters.Category,
		Url:      req.Filters.Url,
		Type:     globals.FileTypeFile, // We always want to agg on just files (data), no directories
		Aggs:     req.Agg,
	}
	query, err := eQuery.AggsQuery()
	if err != nil {
		return nil, err
	}

	out, err := e.conn.Search(req.Filters.Index, req.Filters.Type, nil, query)
	if err != nil {
		return nil, err
	}

	b, err := out.Aggregations.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &search_proto.AggregateResponse{
		Result: string(b),
	}, nil
}
