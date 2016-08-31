package main

import (
	"github.com/micro/go-web"
	"net/http"
)

func main() {
	service := web.NewService(web.Name("go.micro.web.ui"))
	service.Handle("/", http.FileServer(http.Dir("app")))

	service.Init()
	service.Run()
}
