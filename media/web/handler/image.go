package handler

import (
	"encoding/json"
	"fmt"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/custom/proto/custom"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
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
	index := r.FormValue("index")
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

	userId, err := utils.ParseJWTToken(token)
	if err != nil {
		http.Error(w, errors.ErrInvalidUserInCtx.Error(), http.StatusInternalServerError)
		return
	}

	// Build context
	ctx := context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(userId),
	)

	ctx = metadata.NewContext(ctx, map[string]string{
		"Authorization":  token,
		"X-Kazoup-Token": globals.DB_ACCESS_TOKEN,
	})

	// Get file
	rs, err := operations.Read(ctx, &proto_operations.ReadRequest{
		Index: index,
		Type:  globals.FileType,
		Id:    file_id,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("operations.Read %s", err.Error()), http.StatusInternalServerError)
		return
	}

	f, err := file.NewFileFromString(rs.Result)
	if err != nil {
		http.Error(w, "NewFileFromString", http.StatusInternalServerError)
		return
	}

	var fSys fs.Fs
	fSys = ih.getFs(f)

	// Datasource is not on memory yet, (was created after the srv started to run)
	// Lets reload the datasources in memory
	if fSys == nil {
		ih.loadDatasources(ctx, client.DefaultClient)
		fSys = ih.getFs(f)
	}

	// Authorize file system
	auth, err := fSys.Authorize()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update token in DB
	if err := fs.UpdateFsAuth(ctx, fSys.GetDatasourceId(), auth); err != nil {
		log.Println("ERROR", err.Error())
	}

	url := fmt.Sprintf("%s", f.PreviewURL(width, height, mode, quality))

	switch f.GetFileType() {
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
	case globals.GoogleDrive, globals.OneDrive, globals.Dropbox:
		url, err = fSys.GetThumbnail(url)
		if err != nil {
			log.Println("ERROR", err.Error())
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	case globals.Gmail:
		// This is not used in frontend due to tag not able to render base 64 when comes from network request
		b := []byte(url)
		_, err := w.Write(b)
		if err != nil {
			log.Println("ERROR", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return
}

func (ih *ImageHandler) loadDatasources(ctx context.Context, c client.Client) {
	rsp, err := custom.ScrollDatasources(ctx, &proto_custom.ScrollDatasourcesRequest{})
	if err != nil {
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
