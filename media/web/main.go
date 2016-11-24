package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/media/web/handler"
	microweb "github.com/micro/go-web"
)

func main() {
	wd, _ := os.Getwd()

	log.Printf("volume name: %s  path :%s", filepath.VolumeName(wd), wd)

	service := microweb.NewService(microweb.Name("com.kazoup.web.media"))

	service.Handle("/preview", handler.NewImageHandler())
	service.Handle("/thumbnail", handler.NewThumbnailHandler())

	service.Init()
	service.Run()
}
