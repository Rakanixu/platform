package handler

import (
	"encoding/json"
	"fmt"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/fs"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"

	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
)

type ImageHandler struct {
	dbclient         db.DBClient
	datasourceClient datasource.DataSourceClient
	fs               []fs.Fs
}

func NewImageHandler() *ImageHandler {
	ih := &ImageHandler{
		dbclient:         db.NewDBClient(globals.DB_SERVICE_NAME, client.NewClient()),
		datasourceClient: datasource.NewDataSourceClient(globals.DATASOURCE_SERVICE_NAME, client.NewClient()),
		fs:               make([]fs.Fs, 0),
	}

	ih.loadDatasources()

	return ih
}

//ServeHTTP handles requests depending on file type
//http://ADDRESS:8082/desktop/image?file_id={file_id}&width=300&height=300&mode=fit&quality=50
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
	f, err := file.GetFileByID(file_id, ih.dbclient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fSys fs.Fs
	fSys = ih.getFs(f)
	// Datasource is not on memory yet, (was created after the srv started to run)
	// Lets reload the datasources in memory
	if fSys == nil {
		ih.loadDatasources()
		fSys = ih.getFs(f)
	}

	url := fmt.Sprintf("%s", f.PreviewURL(width, height, mode, quality))

	switch f.GetFileType() {
	case globals.Local:
		http.Redirect(w, r, url, http.StatusSeeOther)
	case globals.Slack:
		// There is an issue when setting headers on a redirect request, the headers are dropped
		// We query for the image directly and attach response to the first request response
		c := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		req.Header.Set("Authorization", fSys.Token())
		rsp, err := c.Do(req)
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		defer rsp.Body.Close()
		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println("ERROR", err.Error())
		}
	case globals.GoogleDrive, globals.OneDrive:
		url, err = fSys.GetThumbnail(f.GetIDFromOriginal())
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}

	return
}

func (ih *ImageHandler) loadDatasources() {
	rsp, err := ih.datasourceClient.Search(context.Background(), &datasource.SearchRequest{
		Index: "datasources",
		Type:  "datasource",
		From:  0,
		Size:  9999,
	})
	if err != nil {
		log.Println("ERROR retrieveing datasources for image server")
		return
	}

	var endpoints []*datasource.Endpoint

	if err := json.Unmarshal([]byte(rsp.Result), &endpoints); err != nil {
		log.Println(err.Error())
	}

	for _, v := range endpoints {
		fsfe, err := fs.NewFsFromEndpoint(v)
		if err != nil {
			log.Println(err.Error())
		}
		ih.fs = append(ih.fs, fsfe)
	}
}

func (ih *ImageHandler) getFs(f file.File) fs.Fs {
	for _, v := range ih.fs {
		if v.GetDatasourceId() == f.GetDatasourceID() {
			return v
		}
	}

	return nil
}
