//go:generate cp -r ../src/build/unbundled html
package main

import (
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

func main() {
	// Serve file system
	http.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir("html"))))
	//SPA routes
	http.HandleFunc("/login/", IndexHandler)
	http.HandleFunc("/onboarding", IndexHandler)
	http.HandleFunc("/u/", IndexHandler)
	http.HandleFunc("/search", IndexHandler)
	http.HandleFunc("/settings", IndexHandler)

	err := http.ListenAndServeTLS(":9090", "/secrets/tls.crt", "/secrets/tls.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//IndexHandler serves index.html file of SPA
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/index.html")
}
