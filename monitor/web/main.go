package main

import (
	"fmt"
	"github.com/kardianos/osext"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	proto "github.com/kazoup/platform/monitor/srv/proto/monitor"
	"github.com/kazoup/platform/monitor/web/handler"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"log"
)

func main() {
	dir, err := osext.ExecutableFolder()
	if err != nil {
		log.Println(err.Error())
	}

	service := web.NewService(
		web.Name("com.kazoup.web.monitor"),
		web.Handler(handler.Router()),
	)

	if err := service.Init(); err != nil {
		log.Println(err)
	}

	handler.Init(
		fmt.Sprintf("%s%s", dir, "/templates"),
		proto.NewMonitorClient(globals.MONITOR_SERVICE_NAME, client.DefaultClient),
	)

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
