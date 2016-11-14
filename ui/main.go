package main

import (
	"net/http"

	"github.com/kazoup/platform/lib/globals"
	_ "github.com/micro/go-plugins/transport/tcp"
	"github.com/micro/go-web"
)

func main() {
	service := web.NewService(web.Name(globals.NAMESPACE + ".web.ui"))
	service.Handle("/", http.FileServer(http.Dir("src")))

	service.Init()
	service.Run()
}
