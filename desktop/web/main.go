package main

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-web"
	"net/http"
)

func main() {
	service := web.NewService(web.Name(globals.NAMESPACE + ".web.ui"))
	service.Handle("/", http.FileServer(http.Dir("app")))

	service.Init()
	service.Run()
}
