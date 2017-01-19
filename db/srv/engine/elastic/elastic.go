package elastic

import (
	"encoding/json"
	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/db/srv/engine"
	model "github.com/kazoup/platform/db/srv/engine/elastic/model"
	config "github.com/kazoup/platform/db/srv/proto/config"
	db "github.com/kazoup/platform/db/srv/proto/db"
	subscriber "github.com/kazoup/platform/db/srv/subscriber/elastic"
	"github.com/kazoup/platform/lib/globals"
	search_proto "github.com/kazoup/platform/search/srv/proto/search"
	lib "github.com/mattbaird/elastigo/lib"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"os"
)

type elastic struct {
	*model.Elastic
}

func init() {
	engine.Register(&elastic{
		&model.Elastic{
			FilesChannel:         make(chan *model.FilesChannel),
			SlackUsersChannel:    make(chan *crawler.SlackUserMessage),
			SlackChannelsChannel: make(chan *crawler.SlackChannelMessage),
			CrawlerFinished:      make(chan *crawler.CrawlerFinishedMessage),
		},
	})
}

// Init Elastic db (engine)
// Common init for DB, Config and Subscriber interfaces
func (e *elastic) Init() error {
	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	//connect
	e.Conn = lib.NewConn()
	err := e.Conn.SetFromUrl(url)

	if err != nil {
		return err
	}
	if username != "" {
		e.Conn.Username = username
	}
	if password != "" {
		e.Conn.Password = password
	}
	e.Bulk = e.Conn.NewBulkIndexerErrors(100, 5)
	e.Bulk.BulkMaxDocs = 100000
	e.Bulk.Start()

	// Initialize subscribers
	if err := subscriber.Subscribe(e.Elastic); err != nil {
		return err
	}

	return nil
}

// Create record
func (e *elastic) Create(ctx context.Context, req *db.CreateRequest) (*db.CreateResponse, error) {
	_, err := e.Conn.Index(req.Index, req.Type, req.Id, nil, req.Data)

	return &db.CreateResponse{}, err
}

// Subscribe to crawler file messages
func (e *elastic) SubscribeFiles(ctx context.Context, c client.Client, msg *crawler.FileMessage) error {
	e.FilesChannel <- &model.FilesChannel{
		FileMessage: msg,
		Client:      c,
	}

	return nil
}

// Subscribe to crawler file messages
func (e *elastic) SubscribeSlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error {
	e.SlackUsersChannel <- msg

	return nil
}

// Subscribe to crawler file messages
func (e *elastic) SubscribeSlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error {
	e.SlackChannelsChannel <- msg

	return nil
}

// Subscribe to crawler finished message
func (e *elastic) SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	e.CrawlerFinished <- msg

	return nil
}

func (e *elastic) Read(ctx context.Context, req *db.ReadRequest) (*db.ReadResponse, error) {
	r, err := e.Conn.Get(req.Index, req.Type, req.Id, nil)
	if err != nil {
		// elastigo returns error when does not exists
		if err.Error() == "record not found" {
			return &db.ReadResponse{
				Result: `{}`,
			}, nil
		}

		return &db.ReadResponse{}, err
	}

	data, err := r.Source.MarshalJSON()
	if err != nil {
		return &db.ReadResponse{}, err
	}

	response := &db.ReadResponse{
		Result: string(data),
	}

	return response, nil
}

// Update record
func (e *elastic) Update(ctx context.Context, req *db.UpdateRequest) (*db.UpdateResponse, error) {
	_, err := e.Conn.Index(req.Index, req.Type, req.Id, nil, req.Data)

	return &db.UpdateResponse{}, err
}

// Delete record
func (e *elastic) Delete(ctx context.Context, req *db.DeleteRequest) (*db.DeleteResponse, error) {
	_, err := e.Conn.Delete(req.Index, req.Type, req.Id, nil)

	return &db.DeleteResponse{}, err
}

// DeleteByQuery allows to delete all records that match a DSL query
func (e *elastic) DeleteByQuery(ctx context.Context, req *db.DeleteByQueryRequest) (*db.DeleteByQueryResponse, error) {
	eQuery := ElasticQuery{
		Term:     req.Term,
		Category: req.Category,
		Url:      req.Url,
		Depth:    req.Depth,
		Type:     req.FileType,
		LastSeen: req.LastSeen,
	}

	query, err := eQuery.DeleteQuery()
	if err != nil {
		return nil, err
	}

	_, err = e.Conn.DeleteByQuery(req.Indexes, req.Types, nil, query)
	if err != nil {
		return nil, err
	}

	return &db.DeleteByQueryResponse{}, err
}

// Search ES index
func (e *elastic) Search(ctx context.Context, req *db.SearchRequest) (*db.SearchResponse, error) {
	var results []interface{}
	var err error
	var rstr string
	var uId string

	// Get user id implicitly or explicitly
	if len(req.UserId) == 0 {
		uId, err = globals.ParseJWTToken(ctx)
		if err != nil {
			return &db.SearchResponse{}, err
		}
	} else {
		uId = req.UserId
	}

	eQuery := ElasticQuery{
		Index:                req.Index,
		UserId:               uId,
		Term:                 req.Term,
		From:                 req.From,
		Size:                 req.Size,
		Category:             req.Category,
		Url:                  req.Url,
		Depth:                req.Depth,
		Type:                 req.Type,
		FileType:             req.FileType,
		LastSeen:             req.LastSeen,
		Access:               req.Access,
		NoKazoupFileOriginal: req.NoKazoupFileOriginal,
	}
	query, err := eQuery.Query()

	if err != nil {
		return &db.SearchResponse{}, err
	}

	out, err := e.Conn.Search(req.Index, req.Type, nil, query)
	// Library returns an error when no matching documents, continue
	if err != nil && err.Error() != globals.ES_NO_RESULTS {
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

	if len(results) == 0 {
		rstr = `[]`
	} else {
		b, err := json.Marshal(results)
		if err != nil {
			return &db.SearchResponse{}, err
		}
		rstr = string(b)
	}

	return &db.SearchResponse{
		Result: rstr,
		Info:   info.String(),
	}, nil
}

// SearchById is a way around for read method over aliases
// Since we can't get the record bye id from aliases we need to use search request
// This should return single ID as all files should have unique ID's as we seting them up based on unique path MD5
// Method will work on any index and alias as long ID's are unique
func (e *elastic) SearchById(ctx context.Context, req *db.SearchByIdRequest) (*db.SearchByIdResponse, error) {
	uId, err := globals.ParseJWTToken(ctx)
	if err != nil {
		return &db.SearchByIdResponse{}, err
	}

	eQuery := ElasticQuery{
		Index:  req.Index,
		Id:     req.Id,
		UserId: uId,
		Type:   req.Type,
	}
	query, err := eQuery.QueryById()
	if err != nil {
		return &db.SearchByIdResponse{}, err
	}

	//Search request
	out, err := e.Conn.Search(req.Index, req.Type, nil, query)
	if err != nil {
		return &db.SearchByIdResponse{}, err
	}
	v := out.Hits.Hits
	// hmmm hacky FIXME
	if out.Hits.Total == 0 || out.Hits.Total != 1 {
		return &db.SearchByIdResponse{
			Result: `{}`,
		}, nil
	}

	// Now we should have only one result
	data, err := v[0].Source.MarshalJSON()
	if err != nil {
		return &db.SearchByIdResponse{}, err
	}
	s, err := engine.TypeFactory(req.Type, string(data))
	if err != nil {
		return &db.SearchByIdResponse{}, err
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return &db.SearchByIdResponse{}, err
	}

	return &db.SearchByIdResponse{
		Result: string(data),
	}, nil
}

// CreateIndexWithSettings creates an ES index with settings
func (e *elastic) CreateIndex(ctx context.Context, req *config.CreateIndexRequest) (*config.CreateIndexResponse, error) {
	exists, err := e.Conn.IndicesExists(req.Index)
	if err != nil {
		return &config.CreateIndexResponse{}, err
	}

	if !exists {
		_, err := e.Conn.CreateIndex(req.Index)
		if err != nil {
			return &config.CreateIndexResponse{}, err
		}
	}

	return &config.CreateIndexResponse{}, nil
}

// Status elasticsearch cluster
func (e *elastic) Status(ctx context.Context, req *config.StatusRequest) (*config.StatusResponse, error) {
	clusterState, err := e.Conn.ClusterState(lib.ClusterStateFilter{
		FilterNodes:        true,
		FilterRoutingTable: true,
		FilterMetadata:     true,
		FilterBlocks:       true,
	})

	b, err := json.Marshal(clusterState)
	if err != nil {
		return &config.StatusResponse{}, err
	}

	response := &config.StatusResponse{
		Status: string(b),
	}

	return response, err
}

// AddAlias to assign indexes (aliases) per datasource
func (e *elastic) AddAlias(ctx context.Context, req *config.AddAliasRequest) (*config.AddAliasResponse, error) {
	_, err := e.Conn.AddAlias(req.Index, req.Alias)
	if err != nil {
		return nil, err
	}

	return &config.AddAliasResponse{}, nil
}

// DeleteIndex from ES
func (e *elastic) DeleteIndex(ctx context.Context, req *config.DeleteIndexRequest) (*config.DeleteIndexResponse, error) {
	_, err := e.Conn.DeleteIndex(req.Index)
	if err != nil {
		return nil, err
	}

	return &config.DeleteIndexResponse{}, nil
}

// DeleteAlias from ES
func (e *elastic) DeleteAlias(ctx context.Context, req *config.DeleteAliasRequest) (*config.DeleteAliasResponse, error) {
	_, err := e.RemoveAlias(req.Index, req.Alias)
	if err != nil {
		return nil, err
	}

	return &config.DeleteAliasResponse{}, nil
}

// RenameAlias from ES
func (e *elastic) RenameAlias(ctx context.Context, req *config.RenameAliasRequest) (*config.RenameAliasResponse, error) {
	var err error

	_, err = e.RemoveAlias(req.Index, req.OldAlias)
	if err != nil {
		return nil, err
	}

	_, err = e.AddAlias(ctx, &config.AddAliasRequest{
		Index: req.Index,
		Alias: req.NewAlias,
	})
	if err != nil {
		return nil, err
	}

	return &config.RenameAliasResponse{}, nil
}

// Aggregate allow us to query for aggs in ES
func (e *elastic) Aggregate(ctx context.Context, req *search_proto.AggregateRequest) (*search_proto.AggregateResponse, error) {
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

	out, err := e.Conn.Search(req.Filters.Index, req.Filters.Type, nil, query)
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
