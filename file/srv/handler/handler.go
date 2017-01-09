package handler

import (
	"encoding/json"
	proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

// File struct
type File struct {
	Client client.Client
}

// Create File handler
func (f *File) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.DatasourceId) == 0 {
		return errors.BadRequest("com.kazoup.srv.file.Create", "datasource_id required")
	}

	// Instantiate file system
	fsys, err := NewFileSystem(f.Client, ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := fsys.Authorize()
	if err != nil {
		return err
	}
	log.Println("!!", auth.Expiry)
	// Update token in DB
	if err := UpdateFileSystemAuth(f.Client, ctx, req.DatasourceId, auth); err != nil {
		return err
	}

	// Create a file for given file system
	ch := fsys.Create(*req)
	// Block while creating
	fmc := <-ch
	close(ch)

	// Check for errors that happened in the goroutine
	if fmc.Error != nil {
		return fmc.Error
	}

	// Index created file
	if err := file.IndexAsync(f.Client, fmc.File, globals.FilesTopic, fmc.File.GetIndex(), true); err != nil {
		return err
	}

	b, err := json.Marshal(fmc.File)
	if err != nil {
		return err
	}

	rsp.Data = string(b)
	rsp.DocUrl = fmc.File.GetURL()

	return nil
}

// Delete File handler
func (f *File) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.DatasourceId) == 0 {
		return errors.BadRequest("com.kazoup.srv.file", "datasource_id required")
	}

	// Instantiate file system
	fsys, err := NewFileSystem(f.Client, ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	_, err = fsys.DeleteFile(ctx, f.Client, *req)
	if err != nil {
		return err
	}

	return nil
}

// Share file handler
func (f *File) Share(ctx context.Context, req *proto.ShareRequest, rsp *proto.ShareResponse) error {
	if len(req.OriginalId) == 0 {
		return errors.BadRequest("com.kazoup.srv.file", "file original_id required")
	}

	if len(req.DatasourceId) == 0 {
		return errors.BadRequest("com.kazoup.srv.file", "datasource_id required")
	}

	// Instantiate file system
	fsys, err := NewFileSystem(f.Client, ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	url, err := fsys.ShareFile(ctx, f.Client, *req)
	if err != nil {
		return err
	}

	rsp.PublicUrl = url
	rsp.SharePublicly = req.SharePublicly // Return it back for frontend callback handler

	return nil
}
