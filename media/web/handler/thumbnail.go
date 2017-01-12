package handler

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/globals"
	"net/http"
)

type ThumbnailHandler struct{}

func NewThumbnailHandler() *ThumbnailHandler {
	return &ThumbnailHandler{}
}

// ServeHTTP redirect request to google cloud storage where thumbnails are stored
func (th *ThumbnailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ncs, err := cs.NewCloudStorageFromEndpoint(&datasource_proto.Endpoint{
		Index: r.FormValue("index"),
	}, globals.GoogleCloudStorage)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR Instantiating cloud storage %s", err), http.StatusBadRequest)
		return
	}

	url, err := ncs.SignedObjectStorageURL(r.FormValue("index"), r.FormValue("file_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR getting signed URL %s", err), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)

	return
}
