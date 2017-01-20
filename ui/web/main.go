//go:generate cp -r ../src/build/unbundled html
package main

import (
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

func main() {
	go http.ListenAndServe(":9091", http.HandlerFunc(redirect))

	// Serve file system
	http.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir("html"))))
	//SPA routes
	http.HandleFunc("/login/", IndexHandler)
	http.HandleFunc("/onboarding", IndexHandler)
	http.HandleFunc("/u/", IndexHandler)
	http.HandleFunc("/search", IndexHandler)
	http.HandleFunc("/settings", IndexHandler)

	err := http.ListenAndServeTLS(":9090", "/secrets/all.pem", "/secrets/tls.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//IndexHandler serves index.html file of SPA
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	http.ServeFile(w, r, "html/index.html")
}

// redirect redirects request over http to https
func redirect(w http.ResponseWriter, req *http.Request) {
	// https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
}
