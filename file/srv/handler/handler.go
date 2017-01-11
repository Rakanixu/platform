package handler

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// File struct
type File struct {
	Client client.Client
}

// Create File handler
func (f *File) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.DatasourceId) == 0 {
		return errors.BadRequest("com.kazoup.srv.file", "datasource_id required")
	}

	// Instantiate file system
	fsys, err := NewFileSystem(f.Client, ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Create a file for given file system
	r, err := fsys.CreateFile(ctx, f.Client, *req)
	if err != nil {
		return err
	}

	rsp.Data = r.Data
	rsp.DocUrl = r.DocUrl

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

	// Delete file from GCS
	if err := f.Client.Publish(globals.NewSystemContext(), f.Client.NewPublication(globals.DeleteFileInBucketTopic, &datasource_proto.DeleteFileInBucketMessage{
		FileId: req.FileId,
		Index:  req.Index,
	})); err != nil {
		fmt.Println("ERROR cleaning thumbnail", err)
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
