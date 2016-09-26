package main

import (
	"log"

	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/kazoup/platform/web/web/handler"
	"github.com/micro/go-web"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name(globals.NAMESPACE+".web.web"),
		web.Version("latest"),
		web.Address("127.0.0.1:8083"),
	)
	c := wrappers.NewKazoupClient()

	// register call handler
	//service.HandleFunc("/rpc", handler.RPC)
	service.Handle("/", &handler.Handler{
		Client: c,
	})

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
