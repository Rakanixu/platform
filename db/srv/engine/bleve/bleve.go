package bleve

import (
	"bytes"
	"encoding/json"
	"errors"
	lib "github.com/blevesearch/bleve"
	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"

	"github.com/kazoup/platform/db/srv/engine"
	db "github.com/kazoup/platform/db/srv/proto/db"
	search_proto "github.com/kazoup/platform/search/srv/proto/search"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"sync"
)

var (
	dataPath        = os.TempDir()
	kazoupNamespace = "/kazoup/"
	filesIndex      = "files"
)

type bleve struct {
	mu              sync.Mutex
	indexMap        map[string]lib.Index
	filesChannel    chan *crawler.FileMessage
	crawlerFinished chan *crawler.CrawlerFinishedMessage
	batchSize       int
}

func init() {
	engine.Register(&bleve{
		indexMap:        make(map[string]lib.Index),
		filesChannel:    make(chan *crawler.FileMessage),
		crawlerFinished: make(chan *crawler.CrawlerFinishedMessage),
		batchSize:       1000,
	})
}

// Init bleve db / indexes
func (b *bleve) Init() error {
	err := errors.New("")
	files, err := ioutil.ReadDir(os.TempDir() + kazoupNamespace)
	// Creates directory if not exists, but err is not nil
	if err != err {
		return err
	}

	b.mu.Lock()
	for _, file := range files {
		if err := openIndex(b, file.Name()); err != nil {
			return err
		}
	}

	// Check files index exists, if not, create. Subscriber needs to have it open
	if b.indexMap[filesIndex] == nil {
		if err := openIndex(b, filesIndex); err != nil {
			return err
		}
	}
	b.mu.Unlock()

	if err := indexer(b); err != nil {
		return err
	}

	if err := enricher(b); err != nil {
		return err
	}

	return nil
}

// Subscribe to crawler file messages
func (b *bleve) SubscribeFiles(ctx context.Context, msg *crawler.FileMessage) error {
	b.filesChannel <- msg

	return nil
}

// Subscribe to crawler file messages
func (b *bleve) SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	b.crawlerFinished <- msg

	return nil
}

// Create record
func (b *bleve) Create(req *db.CreateRequest) (*db.CreateResponse, error) {
	response := &db.CreateResponse{}
	ds := &datasource_proto.Endpoint{}

	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	if req.Type == "datasource" {
		if err := json.Unmarshal([]byte(req.Data), ds); err != nil {
			return response, err
		}
	}

	return response, b.indexMap[req.Index].Index(req.Id, ds)
}

// Read record
func (b *bleve) Read(req *db.ReadRequest) (*db.ReadResponse, error) {
	response := &db.ReadResponse{}

	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	document, err := b.indexMap[req.Index].Document(req.Id)
	if err != nil {
		return response, err
	}

	jsonObj := gabs.New()
	if document != nil {
		for _, v := range document.Fields {
			jsonObj.Set(string(v.Value()), v.Name())
		}
	}
	response.Result = jsonObj.String()

	return response, err
}

// Update record
func (b *bleve) Update(req *db.UpdateRequest) (*db.UpdateResponse, error) {
	response := &db.UpdateResponse{}

	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	return response, b.indexMap[req.Index].Index(req.Id, req.Data)
}

// Delete record
func (b *bleve) Delete(req *db.DeleteRequest) (*db.DeleteResponse, error) {
	response := &db.DeleteResponse{}

	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	return response, b.indexMap[req.Index].Delete(req.Id)
}

// CreateIndexWithSettings creates an bleve index
func (b *bleve) CreateIndexWithSettings(req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error) {
	response := &db.CreateIndexWithSettingsResponse{}

	if b.indexExists(req.Index) {
		return response, errors.New("Index already exists")
	}

	if err := openIndex(b, req.Index); err != nil {
		return response, err
	}

	return response, nil
}

// PutMappingFromJSON not required for bleve
func (b *bleve) PutMappingFromJSON(req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error) {
	response := &db.PutMappingFromJSONResponse{}

	return response, nil
}

// Status bleve indexes
func (b *bleve) Status(req *db.StatusRequest) (*db.StatusResponse, error) {
	response := &db.StatusResponse{}

	jsonStatus := gabs.New()
	jsonStatus.SetP("bleve", "master_node")
	jsonStatus.SetP(nil, "metadata.indices")

	for k, _ := range b.indexMap {
		if b.indexExists(k) {
			jsonStatus.SetP("open", "metadata.indices."+k+".state")

		}
	}

	response.Status = jsonStatus.String()

	return response, nil
}

func (b *bleve) Search(req *db.SearchRequest) (*db.SearchResponse, error) {
	var indexSearch string
	var sr *lib.SearchRequest
	var buffer bytes.Buffer
	count := 0

	queries := []lib.Query{}
	prefixQueries := []lib.Query{}

	if len(req.Index) > 0 {
		indexSearch = req.Index
	} else {
		indexSearch = filesIndex
	}

	if b.indexMap[indexSearch] == nil {
		return &db.SearchResponse{}, errors.New("index does not exists")
	}

	if len(req.Term) > 0 {
		termQuery := lib.NewMatchQuery(req.Term)
		queries = append(queries, termQuery)
	}
	if len(req.Category) > 0 {
		categoryQuery := lib.NewTermQuery(req.Category)
		categoryQuery.SetField("category")
		queries = append(queries, categoryQuery)
	}
	if req.Depth > 0 {
		min := new(float64)
		max := new(float64)
		*min = float64(req.Depth)
		*max = float64(req.Depth + 1)
		depthQuery := lib.NewNumericRangeQuery(min, max)
		depthQuery.SetField("depth")
		queries = append(queries, depthQuery)
	}
	if len(req.Url) > 0 {
		urlQuery := lib.NewPrefixQuery(req.Url)
		urlQuery.SetField("url")
		prefixQueries = append(prefixQueries, urlQuery)
	}

	if indexSearch != filesIndex {
		allQuery := lib.NewMatchAllQuery()
		queries = append(queries, allQuery)
	}

	query := lib.NewConjunctionQuery([]lib.Query{
		prefixQueries[0],
		lib.NewConjunctionQuery(queries),
	})

	sr = lib.NewSearchRequestOptions(query, int(req.Size), int(req.From), false)
	sr.Fields = []string{"*"} // Retrieve all fields

	results, err := b.indexMap[indexSearch].Search(sr)
	if err != nil {
		return &db.SearchResponse{}, err
	}

	buffer.WriteString(`[`)
	for _, obj := range results.Hits {
		file, err := json.Marshal(obj.Fields)
		if err != nil {
			return &db.SearchResponse{}, err
		}
		buffer.WriteString(string(file))

		if count < len(results.Hits)-1 {
			buffer.WriteString(`,`)
		}
		count++
	}
	buffer.WriteString(`]`)

	jsonInfo := gabs.New()
	jsonInfo.SetP(results.Total, "total")

	return &db.SearchResponse{
		Result: buffer.String(),
		Info:   jsonInfo.String(),
	}, nil
}

// RenameIndexAlias to assign indexes (aliases) per datasource
func (b *bleve) AddAlias(req *db.AddAliasRequest) (*db.AddAliasResponse, error) {
	// TODO: implement
	return &db.AddAliasResponse{}, nil
}

// DeleteIndex from bleve
func (b *bleve) DeleteIndex(req *db.DeleteIndexRequest) (*db.DeleteIndexResponse, error) {
	// TODO: implement
	return &db.DeleteIndexResponse{}, nil
}

// DeleteAlias from bleve
func (b *bleve) DeleteAlias(req *db.DeleteAliasRequest) (*db.DeleteAliasResponse, error) {
	// TODO: implement
	return &db.DeleteAliasResponse{}, nil
}

// RenameAlias from bleve
func (b *bleve) RenameAlias(req *db.RenameAliasRequest) (*db.RenameAliasResponse, error) {
	// TODO: implement
	return &db.RenameAliasResponse{}, nil
}

func (b *bleve) Aggregate(req *search_proto.AggregateRequest) (*search_proto.AggregateResponse, error) {
	return nil, nil
}

func (b *bleve) DeleteByQuery(req *db.DeleteByQueryRequest) (*db.DeleteByQueryResponse, error) {
	return nil, nil
}
func (b *bleve) SearchById(req *db.SearchByIdRequest) (*db.SearchByIdResponse, error) {
	return nil, nil
}

func (b *bleve) SubscribeSlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error {
	return nil
}
func (b *bleve) SubscribeSlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error {
	return nil
}
