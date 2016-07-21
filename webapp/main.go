package main

import (
	"log"
	"net/http"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-web"
)

func main() {

	var ServeFrom string = "web/build/bundled"
	cmd.Init()
	// New login service
	service := web.NewService(
		web.Name("go.micro.web.webapp"),
		web.Version("0.0.1"),
	)

	// Init and run service
	if err := service.Init(); err != nil {
		log.Fatalf("%v", err)
	}
	// Service handlers
	service.Handle(
		"/",
		http.FileServer(http.Dir(ServeFrom)),
	)

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
