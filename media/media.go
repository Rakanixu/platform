package media

import (
	"github.com/kazoup/platform/media/web/handler"
	"github.com/micro/cli"
	microweb "github.com/micro/go-web"
	"log"
	"os"
	"path/filepath"
)

func web(ctx *cli.Context) {
	wd, _ := os.Getwd()

	log.Printf("volume name: %s  path :%s", filepath.VolumeName(wd), wd)

	service := microweb.NewService(microweb.Name("go.micro.web.media"))

	service.Handle("/preview", handler.NewImageHandler())
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
