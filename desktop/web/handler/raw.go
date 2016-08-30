package handler

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"path"
)

type RAWHandler struct {
	root string
}

func NewRAWHandler(root string) *RAWHandler {
	return &RAWHandler{root}
}

func (s *RAWHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Join(s.root, r.URL.Path)
	w.Header()["Content-Type"] = []string{"video/mp4"}
	w.Header()["Transfer-Encoding"] = []string{"chunked"}
	w.Header()["Accept-Ranges"] = []string{"bytes"}
	log.Printf("Headers : %s", w.Header())

	f, err := os.Open(path)
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	b := bufio.NewReader(f)

	if _, err := b.WriteTo(w); err != nil {
		log.Print(err)
	}

}
