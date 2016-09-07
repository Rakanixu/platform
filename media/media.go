package media

import (
	"crypto/sha256"
	"github.com/kazoup/platform/media/web/handler"
	"github.com/micro/cli"
	microweb "github.com/micro/go-web"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_gift "github.com/pierrre/imageserver/http/gift"
	imageserver_http_image "github.com/pierrre/imageserver/http/image"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/gif"
	imageserver_image_gift "github.com/pierrre/imageserver/image/gift"
	_ "github.com/pierrre/imageserver/image/jpeg"
	_ "github.com/pierrre/imageserver/image/png"
	imageserver_file "github.com/pierrre/imageserver/source/file"
	"golang.org/x/net/websocket"

	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
)

var (
	flagMemory = int64(128 * (1 << 20))
)

func web(ctx *cli.Context) {

	wd, _ := os.Getwd()
	contentDir := "/"
	log.Printf("volume name: %s  path :%s", filepath.VolumeName(wd), wd)
	service := microweb.NewService(microweb.Name("go.micro.web.media"))
	//http://localhost:8082/desktop/image?source={file_id}&width=300&height=300&mode=fit&quality=50
	service.Handle("/image", handler.NewImageHandler())
	service.Handle("/image/local", &imageserver_http.Handler{
		Parser: imageserver_http.ListParser([]imageserver_http.Parser{
			&imageserver_http.SourceParser{},
			&imageserver_http_gift.ResizeParser{},
			&imageserver_http_image.FormatParser{},
			&imageserver_http_image.QualityParser{},
		}),
		Server: &imageserver.HandlerServer{
			Server: newServerMemory(&imageserver_file.Server{}),
			Handler: &imageserver_image.Handler{
				Processor: &imageserver_image_gift.ResizeProcessor{},
			},
		},
	})
	service.Handle("/stream/", http.StripPrefix("/stream/", handler.NewPlaylistHandler(contentDir)))
	service.Handle("/frame/", http.StripPrefix("/frame/", handler.NewFrameHandler(contentDir)))
	service.Handle("/segments/", http.StripPrefix("/segments/", handler.NewStreamHandler(contentDir)))
	service.Handle("/mp4/", http.StripPrefix("/mp4/", handler.NewMP4Handler(contentDir)))
	service.Handle("/raw/", http.StripPrefix("/raw/", handler.NewRAWHandler(contentDir)))
	service.Handle("/webm/", http.StripPrefix("/webm/", handler.NewWebmHandler(contentDir)))

	//TODO move to crawler web service
	service.Handle("/crawler/status", websocket.Handler(handler.CrawlerStatus))

	service.Run()
}

func newServerMemory(srv imageserver.Server) imageserver.Server {
	if flagMemory <= 0 {
		return srv
	}
	cch := imageserver_cache_memory.New(flagMemory)
	kg := imageserver_cache.NewParamsHashKeyGenerator(sha256.New)
	return &imageserver_cache.Server{
		Server:       srv,
		Cache:        cch,
		KeyGenerator: kg,
	}
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
