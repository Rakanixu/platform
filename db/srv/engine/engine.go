package engine

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db "github.com/kazoup/platform/db/srv/proto/db"
	flag_proto "github.com/kazoup/platform/flag/srv/proto/flag"
	"github.com/kazoup/platform/lib/file"
	search_proto "github.com/kazoup/platform/search/srv/proto/search"
	"golang.org/x/net/context"
)

const (
	File       = "file"
	Datasource = "datasource"
	Flag       = "flag"
)

type Engine interface {
	Init() error
	SubscribeFiles(ctx context.Context, msg *crawler.FileMessage) error
	SubscribeSlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error
	SubscribeSlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error
	SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error
	Create(ctx context.Context, req *db.CreateRequest) (*db.CreateResponse, error)
	Read(ctx context.Context, req *db.ReadRequest) (*db.ReadResponse, error)
	Update(ctx context.Context, req *db.UpdateRequest) (*db.UpdateResponse, error)
	Delete(ctx context.Context, req *db.DeleteRequest) (*db.DeleteResponse, error)
	DeleteByQuery(ctx context.Context, req *db.DeleteByQueryRequest) (*db.DeleteByQueryResponse, error)
	CreateIndexWithSettings(ctx context.Context, req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error)
	PutMappingFromJSON(ctx context.Context, req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error)
	Status(ctx context.Context, req *db.StatusRequest) (*db.StatusResponse, error)
	Search(ctx context.Context, req *db.SearchRequest) (*db.SearchResponse, error)
	SearchById(ctx context.Context, req *db.SearchByIdRequest) (*db.SearchByIdResponse, error)
	AddAlias(ctx context.Context, req *db.AddAliasRequest) (*db.AddAliasResponse, error)
	DeleteIndex(ctx context.Context, req *db.DeleteIndexRequest) (*db.DeleteIndexResponse, error)
	DeleteAlias(ctx context.Context, req *db.DeleteAliasRequest) (*db.DeleteAliasResponse, error)
	RenameAlias(ctx context.Context, req *db.RenameAliasRequest) (*db.RenameAliasResponse, error)
	Aggregate(ctx context.Context, req *search_proto.AggregateRequest) (*search_proto.AggregateResponse, error)
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

func SubscribeFiles(ctx context.Context, msg *crawler.FileMessage) error {
	return engine.SubscribeFiles(ctx, msg)
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

func CreateIndexWithSettings(ctx context.Context, req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error) {
	return engine.CreateIndexWithSettings(ctx, req)
}

func PutMappingFromJSON(ctx context.Context, req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error) {
	return engine.PutMappingFromJSON(ctx, req)
}

func Status(ctx context.Context, req *db.StatusRequest) (*db.StatusResponse, error) {
	return engine.Status(ctx, req)
}

func Search(ctx context.Context, req *db.SearchRequest) (*db.SearchResponse, error) {
	return engine.Search(ctx, req)
}

func SearchById(ctx context.Context, req *db.SearchByIdRequest) (*db.SearchByIdResponse, error) {
	return engine.SearchById(ctx, req)
}

func AddAlias(ctx context.Context, req *db.AddAliasRequest) (*db.AddAliasResponse, error) {
	return engine.AddAlias(ctx, req)
}

func DeleteIndex(ctx context.Context, req *db.DeleteIndexRequest) (*db.DeleteIndexResponse, error) {
	return engine.DeleteIndex(ctx, req)
}

func DeleteAlias(ctx context.Context, req *db.DeleteAliasRequest) (*db.DeleteAliasResponse, error) {
	return engine.DeleteAlias(ctx, req)
}

func RenameAlias(ctx context.Context, req *db.RenameAliasRequest) (*db.RenameAliasResponse, error) {
	return engine.RenameAlias(ctx, req)
}

func Aggregate(ctx context.Context, req *search_proto.AggregateRequest) (*search_proto.AggregateResponse, error) {
	return engine.Aggregate(ctx, req)
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
	case Flag:
		return &flag_proto.ReadResponse{}, nil
	}

	return nil, nil
}
