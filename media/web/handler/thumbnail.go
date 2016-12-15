package handler

import (
	"fmt"
	"github.com/kazoup/platform/lib/fs"
	"net/http"
)

type ThumbnailHandler struct{}

func NewThumbnailHandler() *ThumbnailHandler {
	return &ThumbnailHandler{}
}

// ServeHTTP redirect request to google cloud storage where thumbnails are stored
func (th *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url, err := fs.SignedObjectStorageURL(r.FormValue("index"), r.FormValue("file_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR ThumbnailHandler %s", err), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)

	return
}
