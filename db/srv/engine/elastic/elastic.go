package elastic

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/db/srv/engine"
	model "github.com/kazoup/platform/db/srv/engine/elastic/model"
	config "github.com/kazoup/platform/db/srv/proto/config"
	db "github.com/kazoup/platform/db/srv/proto/db"
	subscriber "github.com/kazoup/platform/db/srv/subscriber/elastic"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
	"log"
	"os"
	"time"
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
func (e *elastic) Init(c client.Client) error {
	var err error

	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	// Client
	e.Client, err = elib.NewSimpleClient(
		elib.SetURL(url),
		elib.SetBasicAuth(username, password),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		return err
	}

	rs, err := globals.NewUUID()
	if err != nil {
		return err
	}

	// Bulk Processor, used for users and channels
	e.BulkProcessor, err = e.Client.BulkProcessor().
		Name(fmt.Sprintf("bulkProcessor-%s", rs)).
		Workers(3).
		BulkActions(100).                // commit if # requests >= 100
		BulkSize(2 << 20).               // commit if size of requests >= 2 MB, probably to big, btw other constrains will be hit before
		FlushInterval(10 * time.Second). // commit every 10s
		Do(context.Background())
	if err != nil {
		return err
	}

	// Bulk Files Processor, used for index After function to
	e.BulkFilesProcessor, err = e.Client.BulkProcessor().
		After(func(executionId int64, requests []elib.BulkableRequest, response *elib.BulkResponse, err error) {
			for _, req := range requests {
				type updateBody struct {
					Doc *file.KazoupFile `json:"doc"`
				}

				var kf updateBody

				// elib.BulkableRequest stores two objects, headers and body
				src, err := req.Source()
				if err != nil {
					log.Println("Error: %v", err)
					return
				}

				if len(src) == 2 {
					json.Unmarshal([]byte(src[1]), &kf)
				}

				n := &enrich_proto.EnrichMessage{
					Index:  kf.Doc.Index,
					Id:     kf.Doc.ID,
					UserId: kf.Doc.UserId,
				}

				if err := c.Publish(globals.NewSystemContext(), c.NewPublication(globals.ThumbnailTopic, n)); err != nil {
					log.Print("Publishing ThumbnailTopic error %s", err)
				}

				var topic string
				publishMsg := true

				switch kf.Doc.Category {
				case globals.CATEGORY_DOCUMENT:
					topic = globals.DocEnrichTopic
				case globals.CATEGORY_PICTURE:
					topic = globals.ImgEnrichTopic
				case globals.CATEGORY_AUDIO:
					topic = globals.AudioEnrichTopic
				default:
					publishMsg = false
				}

				// Publish EnrichMessage
				if publishMsg {
					if err := c.Publish(globals.NewSystemContext(), c.NewPublication(topic, n)); err != nil {
						log.Print("Publishing (enrich file) error %s", err)
					}

					log.Println("ENRICH MSG SENT", topic, n.Id)
					time.Sleep(globals.PUBLISHING_DELAY_MS)
				}
			}
		}).
		Name(fmt.Sprintf("bulkFilesProcessor-%s", rs)).
		Workers(3).
		BulkActions(100).               // commit if # requests >= 100
		BulkSize(2 << 20).              // commit if size of requests >= 2 MB, probably to big, btw other constrains will be hit before
		FlushInterval(5 * time.Second). // commit every 5s, notification message can be send and until 5s later is not really finished
		Do(context.Background())
	if err != nil {
		return err
	}

	// Initialize subscribers
	if err := subscriber.Subscribe(e.Elastic); err != nil {
		return err
	}

	return nil
}

// Create record
func (e *elastic) Create(ctx context.Context, req *db.CreateRequest) (*db.CreateResponse, error) {
	exists, err := e.Client.IndexExists(req.Index).Do(ctx)
	if err != nil {
		return &db.CreateResponse{}, err
	}

	if !exists {
		// Create a new index.
		_, err := e.Client.CreateIndex(req.Index).Do(ctx)
		if err != nil {
			return &db.CreateResponse{}, err
		}
	}

	_, err = e.Client.Index().Index(req.Index).Type(req.Type).Id(req.Id).BodyString(req.Data).Do(ctx)
	if err != nil {
		return &db.CreateResponse{}, err
	}

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
	r, err := e.Client.Get().Index(req.Index).Type(req.Type).Id(req.Id).Do(ctx)
	if err != nil {
		return &db.ReadResponse{}, err
	}

	// Return empty if not found
	if !r.Found {
		return &db.ReadResponse{}, nil
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
	// Udate library from library is meant to do updates when data is known
	_, err := e.Create(ctx, &db.CreateRequest{
		Index: req.Index,
		Type:  req.Type,
		Id:    req.Id,
		Data:  req.Data,
	})

	return &db.UpdateResponse{}, err
}

// Delete record
func (e *elastic) Delete(ctx context.Context, req *db.DeleteRequest) (*db.DeleteResponse, error) {
	_, err := e.Client.Delete().Index(req.Index).Type(req.Type).Id(req.Id).Do(ctx)

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

	_, err = e.Client.DeleteByQuery(req.Indexes...).Type(req.Types...).Body(query).Do(ctx)
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

	// Get user id from context
	uId, err = globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return nil, err
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
		ContentCategory:      req.ContentCategory,
		NoKazoupFileOriginal: req.NoKazoupFileOriginal,
	}
	query, err := eQuery.Query()
	if err != nil {
		return &db.SearchResponse{}, err
	}

	out, err := e.Client.Search(req.Index).Type(req.Type).Source(query).Do(ctx)
	if err != nil {
		// Error Index does not exists likely to happen (User does not have datasources).
		// Just empty result, as I do not want to check every possible error
		// and manage them properly, empty result or error..
		return &db.SearchResponse{
			Result: `[]`,
			Info:   `{"total":0}`,
		}, nil
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

		// Check if interface implements file.File interface
		_, ok := s.(file.File)
		// Set the highlight
		if ok && v.Highlight != nil && v.Highlight["content"] != nil {
			// We want just one fragment over content field, check how query is generated
			s.(file.File).SetHighlight(v.Highlight["content"][0])
		}

		if err := json.Unmarshal(data, &s); err != nil {
			return &db.SearchResponse{}, err
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
	var uId string
	var err error

	// Get user id implicitly
	uId, err = globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return nil, err
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
	out, err := e.Client.Search(req.Index).Type(req.Type).Source(query).Do(ctx)
	if err != nil {
		return &db.SearchByIdResponse{}, err
	}
	v := out.Hits.Hits
	// hmmm hacky FIXME
	if out.Hits.TotalHits == 0 || out.Hits.TotalHits != 1 {
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

// CreateIndexWithSettings creates an ES index with template settings if match naming
func (e *elastic) CreateIndex(ctx context.Context, req *config.CreateIndexRequest) (*config.CreateIndexResponse, error) {
	exists, err := e.Client.IndexExists(req.Index).Do(ctx)
	if err != nil {
		return &config.CreateIndexResponse{}, err
	}

	if !exists {
		// Create a new index.
		_, err := e.Client.CreateIndex(req.Index).Do(ctx)
		if err != nil {
			return &config.CreateIndexResponse{}, err
		}
	}

	return &config.CreateIndexResponse{}, nil
}

// Status elasticsearch cluster
func (e *elastic) Status(ctx context.Context, req *config.StatusRequest) (*config.StatusResponse, error) {
	cs, err := e.Client.ClusterState().Do(ctx)
	if err != nil {
		return &config.StatusResponse{}, err
	}

	b, err := json.Marshal(cs)
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
	_, err := e.Client.Alias().Add(req.Index, req.Alias).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &config.AddAliasResponse{}, nil
}

// DeleteIndex from ES
func (e *elastic) DeleteIndex(ctx context.Context, req *config.DeleteIndexRequest) (*config.DeleteIndexResponse, error) {
	_, err := e.Client.DeleteIndex(req.Index).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &config.DeleteIndexResponse{}, nil
}

// DeleteAlias from ES
func (e *elastic) DeleteAlias(ctx context.Context, req *config.DeleteAliasRequest) (*config.DeleteAliasResponse, error) {
	_, err := e.Client.Alias().Remove(req.Index, req.Alias).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &config.DeleteAliasResponse{}, nil
}

// RenameAlias from ES
func (e *elastic) RenameAlias(ctx context.Context, req *config.RenameAliasRequest) (*config.RenameAliasResponse, error) {
	_, err := e.Client.Alias().Remove(req.Index, req.OldAlias).Do(ctx)
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
