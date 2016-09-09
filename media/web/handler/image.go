package handler

import (
	"net/http"

	db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/file"
	"github.com/micro/go-micro/client"
)

type ImageHandler struct {
	dbclient db.DBClient
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		dbclient: db.NewDBClient("", client.NewClient()),
	}
}

//ServeHTTP handles requests depending on file type
//http://localhost:8082/desktop/image?file_id={file_id}&width=300&height=300&mode=fit&quality=50
func (ih *ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Extract values from URL
	file_id := r.FormValue("file_id")
	width := r.FormValue("width")
	height := r.FormValue("height")
	mode := r.FormValue("mode")
	quality := r.FormValue("quality")
	//Handle empty values
	if file_id == "" {
		http.Error(w, "file_id argument in URL can not be empty", http.StatusBadRequest)
		return
	}
	if width == "" {
		width = "300"
	}
	if height == "" {
		height = "300"
	}
	if mode == "" {
		mode = "fit"
	}
	if quality == "" {
		quality = "50"
	}
	f, err := file.GetFileByID(file_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	r.Form.Add("source", f.PreviewURL())
	http.Redirect(w, r, "image/preview", http.StatusSeeOther)

}
