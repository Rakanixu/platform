package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/websocket"

	auth_handler "github.com/kazoup/platform/auth/web/handler"
	"github.com/kazoup/platform/desktop/web/handler"
	"github.com/micro/go-web"
	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_gift "github.com/pierrre/imageserver/http/gift"
	imageserver_http_image "github.com/pierrre/imageserver/http/image"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/gif"
	imageserver_image_gift "github.com/pierrre/imageserver/image/gift"
	_ "github.com/pierrre/imageserver/image/jpeg"
	_ "github.com/pierrre/imageserver/image/png"
	imageserver_file "github.com/pierrre/imageserver/source/file"
)

func main() {
	wd, _ := os.Getwd()
	log.Printf("volume name: %s  path :%s", filepath.VolumeName(wd), wd)
	service := web.NewService(web.Name("go.micro.web.desktop"))
	service.Handle("/", http.FileServer(http.Dir("app")))
	service.Handle("/status", websocket.Handler(handler.Status))
	service.HandleFunc("/google/login", auth_handler.HandleGoogleLogin)
	service.HandleFunc("/GoogleCallback", auth_handler.HandleGoogleCallback)
	service.HandleFunc("/microsoft/login", auth_handler.HandleMicrosoftLogin)
	service.HandleFunc("/microsoft/callback", auth_handler.HandleMicrosoftCallback)
	service.HandleFunc("/slack/login", auth_handler.HandleSlackLogin)
	service.HandleFunc("/slack/callback", auth_handler.HandleSlackCallback)
	//http://localhost:8082/desktop/image?source=/Users/radekdymacz/Pictures/city-wallpaper.jpg&width=300&height=300&mode=fit&quality=50
	service.Handle("/image", &imageserver_http.Handler{
		Parser: imageserver_http.ListParser([]imageserver_http.Parser{
			&imageserver_http.SourceParser{},
			&imageserver_http_gift.ResizeParser{},
			&imageserver_http_image.FormatParser{},
			&imageserver_http_image.QualityParser{},
		}),
		Server: &imageserver.HandlerServer{
			Server: &imageserver_file.Server{},
			Handler: &imageserver_image.Handler{
				Processor: &imageserver_image_gift.ResizeProcessor{},
			},
		},
	})
	service.Handle("/stream/", http.StripPrefix("/stream/", handler.NewPlaylistHandler(contentDir)))
	service.Handle("/frame/", http.StripPrefix("/frame/", handler.NewFrameHandler(contentDir)))
	service.Handle("/segments/", http.StripPrefix("/segments/", handler.NewStreamHandler(contentDir)))

	service.Handle("/webm/", http.StripPrefix("/webm/", handler.NewWebmHandler(contentDir)))
	service.Init()
	service.Run()
}
