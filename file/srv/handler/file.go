package handler

import (
	proto "github.com/kazoup/platform/file/srv/proto/file"
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
	s, err := fsys.CreateFile(req.MimeType)
	if err != nil {
		return err
	}

	rsp.DocUrl = s

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

	url, err := fsys.ShareFile(req.OriginalId)
	if err != nil {
		return err
	}

	rsp.PublicUrl = url

	return nil
}
