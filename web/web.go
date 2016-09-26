package web

import (
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/kazoup/platform/web/web/handler"
	"github.com/micro/cli"
	"log"

	go_web "github.com/micro/go-web"
)

func web(ctx *cli.Context) {

	// create new web service
	service := go_web.NewService(
		go_web.Name(globals.NAMESPACE+".web.proxy"),
		go_web.Version("latest"),
		go_web.Address(":8083"),
	)

	c := wrappers.NewKazoupClient()

	// register call handler
	//service.HandleFunc("/rpc", handler.RPC)
	service.Handle("/", &handler.Handler{
		Client: c,
	})

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func webCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "web",
			Usage:  "Run sheduler web service",
			Action: web,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "proxy",
			Usage:       "web commands",
			Subcommands: webCommands(),
		},
	}
}
