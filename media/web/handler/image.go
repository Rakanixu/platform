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
	"github.com/micro/go-micro/metadata"

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

	//ih.loadDatasources()

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
	token := r.FormValue("token")
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

	// Build context
	ctx := metadata.NewContext(context.TODO(), map[string]string{
		"Authorization": token,
	})

	f, err := file.GetFileByID(ctx, file_id, ih.dbclient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fSys fs.Fs
	fSys = ih.getFs(f)

	// Datasource is not on memory yet, (was created after the srv started to run)
	// Lets reload the datasources in memory
	if fSys == nil {
		ih.loadDatasources(ctx)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rsp.StatusCode != http.StatusOK {
			log.Println("Error getting file status code ", rsp.StatusCode)
			http.Error(w, "Error getting file status code ", http.StatusInternalServerError)
			return
		}
		b, err := ioutil.ReadAll(rsp.Body)
		defer rsp.Body.Close()
		if err != nil {
			log.Println("ERROR", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println("ERROR", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case globals.GoogleDrive, globals.OneDrive:
		url, err = fSys.GetThumbnail(f.GetIDFromOriginal())
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	case globals.Dropbox:
		url, err = fSys.GetThumbnail(f.GetPathDisplay())
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	case globals.Box:
		url, err = fSys.GetThumbnail(f.GetIDFromOriginal())

		log.Println("URL")
		log.Println(url)
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}

	return
}

func (ih *ImageHandler) loadDatasources(ctx context.Context) {
	rsp, err := ih.datasourceClient.Search(ctx, &datasource.SearchRequest{
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
