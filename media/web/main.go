package main

import (
	"github.com/kazoup/platform/media/web/handler"
	microweb "github.com/micro/go-web"
	"log"
	"os"
	"path/filepath"
)

func main() {
	wd, _ := os.Getwd()

	log.Printf("volume name: %s  path :%s", filepath.VolumeName(wd), wd)

	service := microweb.NewService(microweb.Name("go.micro.web.media"))

	service.Handle("/preview", handler.NewImageHandler())

	service.Init()
	service.Run()
}
