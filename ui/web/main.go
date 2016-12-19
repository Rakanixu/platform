//go:generate cp -r ../src/build/unbundled html
package main

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/micro/go-web"
)

func main() {
	service := web.NewService(
		web.Name(globals.NAMESPACE + ".web.ui"),
	)
	// Serve file system
	service.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir("html"))))
	//SPA routes
	service.HandleFunc("/login/", IndexHandler)
	service.HandleFunc("/onboarding", IndexHandler)
	service.HandleFunc("/u/", IndexHandler)
	service.HandleFunc("/search", IndexHandler)
	service.HandleFunc("/settings", IndexHandler)
	service.Init()
	service.Run()
}

//IndexHandler serves index.html file of SPA
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/index.html")
}
