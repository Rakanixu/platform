package engine

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	config "github.com/kazoup/platform/db/srv/proto/config"
	db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/file"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

const (
	File       = "file"
	Datasource = "datasource"
)

type Engine interface {
	Init() error
	DB
	Config
	Subscriber
}

type DB interface {
	Create(ctx context.Context, req *db.CreateRequest) (*db.CreateResponse, error)
	Read(ctx context.Context, req *db.ReadRequest) (*db.ReadResponse, error)
	Update(ctx context.Context, req *db.UpdateRequest) (*db.UpdateResponse, error)
	Delete(ctx context.Context, req *db.DeleteRequest) (*db.DeleteResponse, error)
	DeleteByQuery(ctx context.Context, req *db.DeleteByQueryRequest) (*db.DeleteByQueryResponse, error)
	Search(ctx context.Context, req *db.SearchRequest) (*db.SearchResponse, error)
	SearchById(ctx context.Context, req *db.SearchByIdRequest) (*db.SearchByIdResponse, error)
}

type Config interface {
	CreateIndex(ctx context.Context, req *config.CreateIndexRequest) (*config.CreateIndexResponse, error)
	Status(ctx context.Context, req *config.StatusRequest) (*config.StatusResponse, error)
	AddAlias(ctx context.Context, req *config.AddAliasRequest) (*config.AddAliasResponse, error)
	DeleteIndex(ctx context.Context, req *config.DeleteIndexRequest) (*config.DeleteIndexResponse, error)
	DeleteAlias(ctx context.Context, req *config.DeleteAliasRequest) (*config.DeleteAliasResponse, error)
	RenameAlias(ctx context.Context, req *config.RenameAliasRequest) (*config.RenameAliasResponse, error)
}

type Subscriber interface {
	SubscribeFiles(ctx context.Context, c client.Client, msg *crawler.FileMessage) error
	SubscribeSlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error
	SubscribeSlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error
	SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error
}

var (
	engine Engine
)

func Register(backend Engine) {
	engine = backend
}

func Init() error {
	return engine.Init()
}

type Files struct {
	Client client.Client
}

func (f *Files) SubscribeFiles(ctx context.Context, msg *crawler.FileMessage) error {
	return engine.SubscribeFiles(ctx, f.Client, msg)
}

func SubscribeSlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error {
	return engine.SubscribeSlackUsers(ctx, msg)
}

func SubscribeSlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error {
	return engine.SubscribeSlackChannels(ctx, msg)
}

func SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	return engine.SubscribeCrawlerFinished(ctx, msg)
}

func Create(ctx context.Context, req *db.CreateRequest) (*db.CreateResponse, error) {
	return engine.Create(ctx, req)
}

func Read(ctx context.Context, req *db.ReadRequest) (*db.ReadResponse, error) {
	return engine.Read(ctx, req)
}

func Update(ctx context.Context, req *db.UpdateRequest) (*db.UpdateResponse, error) {
	return engine.Update(ctx, req)
}

func Delete(ctx context.Context, req *db.DeleteRequest) (*db.DeleteResponse, error) {
	return engine.Delete(ctx, req)
}

func DeleteByQuery(ctx context.Context, req *db.DeleteByQueryRequest) (*db.DeleteByQueryResponse, error) {
	return engine.DeleteByQuery(ctx, req)
}

func CreateIndex(ctx context.Context, req *config.CreateIndexRequest) (*config.CreateIndexResponse, error) {
	return engine.CreateIndex(ctx, req)
}

func Status(ctx context.Context, req *config.StatusRequest) (*config.StatusResponse, error) {
	return engine.Status(ctx, req)
}

func Search(ctx context.Context, req *db.SearchRequest) (*db.SearchResponse, error) {
	return engine.Search(ctx, req)
}

func SearchById(ctx context.Context, req *db.SearchByIdRequest) (*db.SearchByIdResponse, error) {
	return engine.SearchById(ctx, req)
}

func AddAlias(ctx context.Context, req *config.AddAliasRequest) (*config.AddAliasResponse, error) {
	return engine.AddAlias(ctx, req)
}

func DeleteIndex(ctx context.Context, req *config.DeleteIndexRequest) (*config.DeleteIndexResponse, error) {
	return engine.DeleteIndex(ctx, req)
}

func DeleteAlias(ctx context.Context, req *config.DeleteAliasRequest) (*config.DeleteAliasResponse, error) {
	return engine.DeleteAlias(ctx, req)
}

func RenameAlias(ctx context.Context, req *config.RenameAliasRequest) (*config.RenameAliasResponse, error) {
	return engine.RenameAlias(ctx, req)
}

func TypeFactory(typ string, data string) (interface{}, error) {
	switch typ {
	case File:
		file, err := file.NewFileFromString(data)
		if err != nil {
			return nil, err
		}
		return file, nil
	case Datasource:
		return &datasource_proto.Endpoint{}, nil
	}

	return nil, nil
}
