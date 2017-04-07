package handler

import (
	"encoding/json"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/file/srv/proto/file"
	db_conn "github.com/kazoup/platform/lib/dbhelper"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type Service struct{}

// Create File handler
func (s *Service) Create(ctx context.Context, req *proto_file.CreateRequest, rsp *proto_file.CreateResponse) error {
	if err := validate.Exists(ctx, req.DatasourceId); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	// Instantiate file system
	fsys, err := NewFileSystem(srv.Client(), ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := fsys.Authorize()
	if err != nil {
		return err
	}

	// Update token in DB
	if err := db_conn.UpdateFileSystemAuth(srv.Client(), ctx, req.DatasourceId, auth); err != nil {
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
	if err := file.IndexAsync(ctx, srv.Client(), fmc.File, globals.FilesTopic, fmc.File.GetIndex(), true); err != nil {
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
func (s *Service) Delete(ctx context.Context, req *proto_file.DeleteRequest, rsp *proto_file.DeleteResponse) error {
	if err := validate.Exists(ctx, req.DatasourceId); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	// Instantiate file system
	fsys, err := NewFileSystem(srv.Client(), ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := fsys.Authorize()
	if err != nil {
		return err
	}

	// Update token in DB
	if err := db_conn.UpdateFileSystemAuth(srv.Client(), ctx, req.DatasourceId, auth); err != nil {
		return err
	}

	ch := fsys.Delete(*req)
	// Block while deleting
	fmc := <-ch
	close(ch)

	// Check for errors that happened in the goroutine
	if fmc.Error != nil {
		return fmc.Error
	}

	// Delete from DB
	_, err = db_conn.DeleteFromDB(srv.Client(), ctx, &db_proto.DeleteRequest{
		Index: req.Index,
		Type:  globals.FileType,
		Id:    req.FileId,
	})
	if err != nil {
		return err
	}

	return nil
}

// Share file handler
func (s *Service) Share(ctx context.Context, req *proto_file.ShareRequest, rsp *proto_file.ShareResponse) error {
	if err := validate.Exists(ctx, req.OriginalId, req.DatasourceId); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	// Instantiate file system
	fsys, err := NewFileSystem(srv.Client(), ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := fsys.Authorize()
	if err != nil {
		return err
	}

	// Update token in DB
	if err := db_conn.UpdateFileSystemAuth(srv.Client(), ctx, req.DatasourceId, auth); err != nil {
		return err
	}

	ch := fsys.Update(*req)
	// Block while updating
	fmc := <-ch
	close(ch)

	// Check for errors that happened in the goroutine
	if fmc.Error != nil {
		return fmc.Error
	}

	b, err := json.Marshal(fmc.File)
	if err != nil {
		return err
	}

	if _, err := db_conn.UpdateFromDB(srv.Client(), ctx, &db_proto.UpdateRequest{
		Index: req.Index,
		Type:  globals.FileType,
		Id:    req.FileId,
		Data:  string(b),
	}); err != nil {
		return err
	}

	rsp.PublicUrl = ""                    // This will change: kbf.Original.SharedLink.URL
	rsp.SharePublicly = req.SharePublicly // Return it back for frontend callback handler

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_file.HealthRequest, rsp *proto_file.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
