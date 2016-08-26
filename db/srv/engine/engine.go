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
	Subscribe(ctx context.Context, msg *crawler.FileMessage) error
	Create(req *db.CreateRequest) (*db.CreateResponse, error)
	Read(req *db.ReadRequest) (*db.ReadResponse, error)
	Update(req *db.UpdateRequest) (*db.UpdateResponse, error)
	Delete(req *db.DeleteRequest) (*db.DeleteResponse, error)
	CreateIndexWithSettings(req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error)
	PutMappingFromJSON(req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error)
	Status(req *db.StatusRequest) (*db.StatusResponse, error)
	Search(req *db.SearchRequest) (*db.SearchResponse, error)
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

func Subscribe(ctx context.Context, msg *crawler.FileMessage) error {
	return engine.Subscribe(ctx, msg)
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
