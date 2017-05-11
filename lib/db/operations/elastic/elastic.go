package elastic

import (
	"encoding/json"
	"github.com/cenkalti/backoff"
	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
	"log"
	"os"
	"time"
)

type elastic struct {
	Client *elib.Client
}

func init() {
	operations.Register(new(elastic))
}

// Init Elastic Operations
func (e *elastic) Init() error {
	var err error

	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	// ElasticSearch Client
	e.Client, err = elib.NewSimpleClient(
		elib.SetURL(url),
		elib.SetBasicAuth(username, password),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		return err
	}

	return nil
}

// Create record
func (e *elastic) Create(ctx context.Context, req *proto_operations.CreateRequest) (*proto_operations.CreateResponse, error) {
	exists, err := e.Client.IndexExists(req.Index).Do(ctx)
	if err != nil {
		return &proto_operations.CreateResponse{}, err
	}

	if !exists {
		// Create a new index.
		_, err := e.Client.CreateIndex(req.Index).Do(ctx)
		if err != nil {
			return &proto_operations.CreateResponse{}, err
		}
	}

	_, err = e.Client.Index().Index(req.Index).Type(req.Type).Id(req.Id).BodyString(req.Data).Do(ctx)
	if err != nil {
		return &proto_operations.CreateResponse{}, err
	}

	return &proto_operations.CreateResponse{}, err
}

// Read record
func (e *elastic) Read(ctx context.Context, req *proto_operations.ReadRequest) (*proto_operations.ReadResponse, error) {
	r, err := e.Client.Get().Index(req.Index).Type(req.Type).Id(req.Id).Do(ctx)
	if err != nil {
		return &proto_operations.ReadResponse{}, err
	}

	// Return empty if not found
	if !r.Found {
		return &proto_operations.ReadResponse{}, nil
	}

	data, err := r.Source.MarshalJSON()
	if err != nil {
		return &proto_operations.ReadResponse{}, err
	}

	response := &proto_operations.ReadResponse{
		Result: string(data),
	}

	return response, nil
}

// Update record
func (e *elastic) Update(ctx context.Context, req *proto_operations.UpdateRequest) (*proto_operations.UpdateResponse, error) {
	var err error
	// FIXME: When working with files, a update can happen in paralel trigerring 409 on ES
	// Backoff will helps us, ES library only accepts interface (as object) and not strings to do parcial updates,
	// so we have to unmarshal the data string into a file, then do the Update
	// We do not have this issue with datasources..
	if req.Type == globals.FileType {
		d, err := operations.TypeFactory(req.Type, req.Data)
		if err != nil {
			return &proto_operations.UpdateResponse{}, err
		}

		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = 200 * time.Millisecond
		bo.MaxInterval = time.Second * 2
		bo.MaxElapsedTime = time.Second * 3

		if err = backoff.Retry(func() error {
			_, err = e.Client.Update().Index(req.Index).Type(req.Type).Id(req.Id).Doc(d).RetryOnConflict(3).Refresh("wait_for").Do(ctx)

			return err
		}, bo); err != nil {
			return &proto_operations.UpdateResponse{}, err
		}
	} else {
		_, err = e.Create(ctx, &proto_operations.CreateRequest{
			Index: req.Index,
			Type:  req.Type,
			Id:    req.Id,
			Data:  req.Data,
		})
	}

	return &proto_operations.UpdateResponse{}, err
}

// Delete record
func (e *elastic) Delete(ctx context.Context, req *proto_operations.DeleteRequest) (*proto_operations.DeleteResponse, error) {
	_, err := e.Client.Delete().Index(req.Index).Type(req.Type).Id(req.Id).Do(ctx)

	return &proto_operations.DeleteResponse{}, err
}

// DeleteByQuery allows to delete all records that match a query
func (e *elastic) DeleteByQuery(ctx context.Context, req *proto_operations.DeleteByQueryRequest) (*proto_operations.DeleteByQueryResponse, error) {
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

	_, err = e.Client.DeleteByQuery(req.Indexes...).Type(req.Types...).Body(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &proto_operations.DeleteByQueryResponse{}, err
}

// Search for records
func (e *elastic) Search(ctx context.Context, req *proto_operations.SearchRequest) (*proto_operations.SearchResponse, error) {
	var results []interface{}
	var err error
	var rstr string
	var uID string

	// Get user id from context
	uID, err = utils.ParseUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	eQuery := ElasticQuery{
		Index:                req.Index,
		UserId:               uID,
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
		ContentCategory:      req.ContentCategory,
		NoKazoupFileOriginal: req.NoKazoupFileOriginal,
	}
	query, err := eQuery.Query()
	if err != nil {
		return &proto_operations.SearchResponse{}, err
	}

	out, err := e.Client.Search(req.Index).Type(req.Type).Source(query).Do(ctx)
	if err != nil {
		// Error Index does not exists likely to happen (User does not have datasources).
		// Just empty result, as I do not want to check every possible error
		// and manage them properly, empty result or error..
		return &proto_operations.SearchResponse{
			Result: `[]`,
			Info:   `{"total":0}`,
		}, nil
	}
	log.Printf("Took : %v ms", out.TookInMillis)
	for _, v := range out.Hits.Hits {
		data, err := v.Source.MarshalJSON()
		if err != nil {
			return &proto_operations.SearchResponse{}, err
		}
		s, err := operations.TypeFactory(req.Type, string(data))
		if err != nil {
			return &proto_operations.SearchResponse{}, err
		}

		// Check if interface implements file.File interface
		_, ok := s.(file.File)
		// Set the highlight
		if ok && v.Highlight != nil && v.Highlight["content"] != nil {
			// We want just one fragment over content field, check how query is generated
			s.(file.File).SetHighlight(v.Highlight["content"][0])
		}

		if err := json.Unmarshal(data, &s); err != nil {
			return &proto_operations.SearchResponse{}, err
		}
		results = append(results, s)
	}

	info := gabs.New()
	info.Set(out.Hits.TotalHits, "total")

	if len(results) == 0 {
		rstr = `[]`
	} else {
		b, err := json.Marshal(results)
		if err != nil {
			return &proto_operations.SearchResponse{}, err
		}
		rstr = string(b)
	}

	return &proto_operations.SearchResponse{
		Result: rstr,
		Info:   info.String(),
	}, nil
}
