package media

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kazoup/platform/media/web/handler"
	"github.com/micro/cli"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	microweb "github.com/micro/go-web"
)

func web(ctx *cli.Context) {
	wd, _ := os.Getwd()

	log.Printf("volume name: %s  path :%s", filepath.VolumeName(wd), wd)

	service := microweb.NewService(microweb.Name("go.micro.web.media"))

	service.Handle("/preview", handler.NewImageHandler())
	service.Handle("/thumbnail", handler.NewThumbnailHandler())

	service.Run()
}

func mediaCommands() []cli.Command {
	return []cli.Command{{
		Name:   "web",
		Usage:  "Run media web service",
		Action: web,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "media",
			Usage:       "Media commands",
			Subcommands: mediaCommands(),
		},
	}
}
