package handler

import (
	"encoding/json"
	"fmt"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_conn "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
)

type ImageHandler struct {
	fs []fs.Fs
}

func NewImageHandler() *ImageHandler {
	ih := &ImageHandler{
		fs: make([]fs.Fs, 0),
	}

	return ih
}

//ServeHTTP handles requests depending on file type
//http://ADDRESS:8082/desktop/image?user_id={user_id}&file_id={file_id}&width=300&height=300&mode=fit&quality=50
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
		"Authorization":  token,
		"X-Kazoup-Token": globals.DB_ACCESS_TOKEN,
	})

	uID, err := globals.ParseJWTToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := file.GetFileByID(ctx, globals.GetMD5Hash(uID), file_id)
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

	// Authorize file system
	auth, err := fSys.Authorize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update token in DB
	if err := db_conn.UpdateFileSystemAuth(client.NewClient(), ctx, fSys.GetDatasourceId(), auth); err != nil {
		log.Println("ERROR", err.Error())
	}

	url := fmt.Sprintf("%s", f.PreviewURL(width, height, mode, quality))

	switch f.GetFileType() {
	case globals.Local:
		http.Redirect(w, r, url, http.StatusSeeOther)
	case globals.Slack, globals.Box:
		// There is an issue when setting headers on a redirect request, the headers are dropped
		// We query for the image directly and attach response to the first request response
		c := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("ERROR", err.Error())
		}

		req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
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
	case globals.Gmail:
		// This is not used in frontend due to tag not able to render base 64 when comes from network request
		b := []byte(`data:image/` + f.GetExtension() + `;base64,` + f.GetBase64())
		_, err := w.Write(b)
		if err != nil {
			log.Println("ERROR", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return
}

func (ih *ImageHandler) loadDatasources(ctx context.Context) {
	//FIXME: mazimun of 9999 datasources. paginate
	req := client.DefaultClient.NewRequest(
		globals.DATASOURCE_SERVICE_NAME,
		"DataSource.Search",
		&datasource.SearchRequest{
			Index: "datasources",
			Type:  "datasource",
			From:  0,
			Size:  9999,
		},
	)
	rsp := &datasource.SearchResponse{}

	if err := client.DefaultClient.Call(ctx, req, rsp); err != nil {
		log.Println("ERROR retrieveing datasources for image server", err)
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
		//FIXME: we are appending to existing ones, so we end up with many copies of same data in memory. just dump and append new data
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
