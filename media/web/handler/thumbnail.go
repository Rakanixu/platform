package handler

import (
	"fmt"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"net/http"
)

type ThumbnailHandler struct{}

func NewThumbnailHandler() *ThumbnailHandler {
	return &ThumbnailHandler{}
}

var gcs *gcslib.GoogleCloudStorage

func init() {
	gcslib.Register()
	gcs = gcslib.NewGoogleCloudStorage()
}

// ServeHTTP redirect request to google cloud storage where thumbnails are stored
func (th *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url, err := gcs.SignedObjectStorageURL(r.FormValue("index"), r.FormValue("file_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR getting signed URL %s", err), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)

	return
}
