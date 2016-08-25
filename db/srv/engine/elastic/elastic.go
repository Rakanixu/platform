package elastic

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/db/srv/engine"
	data "github.com/kazoup/platform/db/srv/engine/elastic/data"
	db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs"
	lib "github.com/mattbaird/elastigo/lib"
	"golang.org/x/net/context"
	"log"
)

type elastic struct {
	conn         *lib.Conn
	bulk         *lib.BulkIndexer
	filesChannel chan *crawler.FileMessage
	esFlags      *[]byte
	esMapping    *[]byte
	esSettings   *[]byte
}

func init() {
	es_flags, err := data.Asset("data/es_flags.json")
	if err != nil {
		// Asset was not found.
		log.Fatal(err)
	}
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
		filesChannel: make(chan *crawler.FileMessage),
		esFlags:      &es_flags,
		esMapping:    &es_mapping,
		esSettings:   &es_settings,
	})
}

func (e *elastic) Init() error {
	e.conn = lib.NewConn()
	e.conn.SetHosts([]string{"localhost:9200"}) //TODO: replace for enterprise version, get flag
	e.bulk = e.conn.NewBulkIndexerErrors(100, 5)
	e.bulk.BulkMaxDocs = 100000
	e.bulk.Start()

	if err := indexer(e); err != nil {
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
func (e *elastic) Subscribe(ctx context.Context, msg *crawler.FileMessage) error {
	e.filesChannel <- msg

	return nil
}

// Read record
func (e *elastic) Read(req *db.ReadRequest) (*db.ReadResponse, error) {
	r, err := e.conn.Get(req.Index, req.Type, req.Id, nil)
	if err != nil {
		return &db.ReadResponse{}, err
	}

	data, err := r.Source.MarshalJSON()
	if err != nil {
		return &db.ReadResponse{}, err
	}

	response := &db.ReadResponse{
		Result: string(data),
	}

	log.Println(response)
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
	var results []*structs.DesktopFile

	eQuery := ElasticQuery{
		Term:     req.Term,
		From:     req.From,
		Size:     req.Size,
		Category: req.Category,
		Url:      req.Url,
		Depth:    req.Depth,
		Type:     req.Type,
	}
	query, err := eQuery.Query()
	if err != nil {
		return &db.SearchResponse{}, nil
	}

	log.Println("query ->", query)
	log.Println("------------------")
	log.Println("req ->", req)

	out, err := e.conn.Search(req.Index, req.Type, nil, query)
	//out, err := lib.Search(req.Index).Type(req.Type). /*.Size(size).From(from)*/ Search(query).Result(e.conn)
	if err != nil {
		return &db.SearchResponse{}, err
	}

	log.Println(out.Hits.Hits)

	for _, v := range out.Hits.Hits {
		var file *structs.DesktopFile
		//log.Println(string(v.Source))
		data, err := v.Source.MarshalJSON()
		if err != nil {
			return &db.SearchResponse{}, err
		}

		if err := json.Unmarshal(data, &file); err != nil {
			return &db.SearchResponse{}, err
		}
		results = append(results, file)
	}

	log.Println(results)

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

	if err := json.Unmarshal(*e.esSettings, &settingsMap); err != nil {
		return &db.CreateIndexWithSettingsResponse{}, err
	}

	_, err := e.conn.CreateIndexWithSettings(req.Index, settingsMap)
	if err != nil {
		return &db.CreateIndexWithSettingsResponse{}, err
	}

	return &db.CreateIndexWithSettingsResponse{}, nil
}

// PutMappingFromJSON puts a mapping into ES
func (e *elastic) PutMappingFromJSON(req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error) {
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
