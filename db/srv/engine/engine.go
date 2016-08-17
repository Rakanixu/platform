package engine

import (
	db "github.com/kazoup/platform/db/srv/proto/db"
)

type Engine interface {
	Init() error
	Create(req *db.CreateRequest) (res *db.CreateResponse, err error)
	Read(req *db.ReadRequest) (res *db.ReadResponse, err error)
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

func Create(req *db.CreateRequest) (res *db.CreateResponse, err error) {
	return engine.Create(req)
}

func Read(req *db.ReadRequest) (res *db.ReadResponse, err error) {
	return engine.Read(req)
}
