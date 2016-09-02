package engine

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db "github.com/kazoup/platform/db/srv/proto/db"
	flag_proto "github.com/kazoup/platform/flag/srv/proto/flag"
	"github.com/kazoup/platform/structs"
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
	SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error
	Create(req *db.CreateRequest) (*db.CreateResponse, error)
	Read(req *db.ReadRequest) (*db.ReadResponse, error)
	Update(req *db.UpdateRequest) (*db.UpdateResponse, error)
	Delete(req *db.DeleteRequest) (*db.DeleteResponse, error)
	CreateIndexWithSettings(req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error)
	PutMappingFromJSON(req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error)
	Status(req *db.StatusRequest) (*db.StatusResponse, error)
	Search(req *db.SearchRequest) (*db.SearchResponse, error)
	AddAlias(req *db.AddAliasRequest) (*db.AddAliasResponse, error)
	DeleteIndex(req *db.DeleteIndexRequest) (*db.DeleteIndexResponse, error)
	DeleteAlias(req *db.DeleteAliasRequest) (*db.DeleteAliasResponse, error)
	RenameAlias(req *db.RenameAliasRequest) (*db.RenameAliasResponse, error)
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

func SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	return engine.SubscribeCrawlerFinished(ctx, msg)
}

func Create(req *db.CreateRequest) (*db.CreateResponse, error) {
	return engine.Create(req)
}

func Read(req *db.ReadRequest) (*db.ReadResponse, error) {
	return engine.Read(req)
}

func Update(req *db.UpdateRequest) (*db.UpdateResponse, error) {
	return engine.Update(req)
}

func Delete(req *db.DeleteRequest) (*db.DeleteResponse, error) {
	return engine.Delete(req)
}

func CreateIndexWithSettings(req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error) {
	return engine.CreateIndexWithSettings(req)
}

func PutMappingFromJSON(req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error) {
	return engine.PutMappingFromJSON(req)
}

func Status(req *db.StatusRequest) (*db.StatusResponse, error) {
	return engine.Status(req)
}

func Search(req *db.SearchRequest) (*db.SearchResponse, error) {
	return engine.Search(req)
}

func AddAlias(req *db.AddAliasRequest) (*db.AddAliasResponse, error) {
	return engine.AddAlias(req)
}

func DeleteIndex(req *db.DeleteIndexRequest) (*db.DeleteIndexResponse, error) {
	return engine.DeleteIndex(req)
}

func DeleteAlias(req *db.DeleteAliasRequest) (*db.DeleteAliasResponse, error) {
	return engine.DeleteAlias(req)
}

func RenameAlias(req *db.RenameAliasRequest) (*db.RenameAliasResponse, error) {
	return engine.RenameAlias(req)
}

func TypeFactory(typ string) interface{} {
	switch typ {
	case File:
		return &structs.DesktopFile{}
	case Datasource:
		return &datasource_proto.Endpoint{}
	case Flag:
		return &flag_proto.ReadResponse{}
	}

	return nil
}
