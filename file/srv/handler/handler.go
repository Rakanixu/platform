package handler

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/db/bulk"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Service struct{}

// Create File handler
func (s *Service) Create(ctx context.Context, req *proto_file.CreateRequest, rsp *proto_file.CreateResponse) error {
	if err := validate.Exists(ctx, req.DatasourceId); err != nil {
		return err
	}

	// Instantiate file system
	fsys, err := NewFileSystem(ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := fsys.Authorize()
	if err != nil {
		return err
	}

	// Update token in DB
	if err := fs.UpdateFsAuth(ctx, req.DatasourceId, auth); err != nil {
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
	b, err := json.Marshal(fmc.File)
	if err != nil {
		return err
	}
	if err := bulk.Files(ctx, &crawler.FileMessage{
		Id:     fmc.File.GetID(),
		UserId: fmc.File.GetUserID(),
		Index:  fmc.File.GetIndex(),
		Notify: false,
		Data:   string(b),
	}); err != nil {
		return err
	}

	rsp.Data = string(b)
	rsp.DocUrl = fmc.File.GetURL()

	return nil
}

// Create File handler
func (s *Service) Read(ctx context.Context, req *proto_file.ReadRequest, rsp *proto_file.ReadResponse) error {
	if err := validate.Exists(ctx, req.Index, req.Id); err != nil {
		return err
	}

	res, err := operations.Read(ctx, &proto_operations.ReadRequest{
		Index: req.Index,
		Type:  globals.FileType,
		Id:    req.Id,
	})
	if err != nil {
		return errors.InternalServerError(globals.FILE_SERVICE_NAME, err.Error())
	}

	rsp.Result = res.Result

	return nil
}

// Delete File handler
func (s *Service) Delete(ctx context.Context, req *proto_file.DeleteRequest, rsp *proto_file.DeleteResponse) error {
	if err := validate.Exists(ctx, req.DatasourceId); err != nil {
		return err
	}

	// Instantiate file system
	fsys, err := NewFileSystem(ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := fsys.Authorize()
	if err != nil {
		return err
	}

	// Update token in DB
	if err := fs.UpdateFsAuth(ctx, req.DatasourceId, auth); err != nil {
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
	_, err = operations.Delete(ctx, &proto_operations.DeleteRequest{
		Index: req.Index,
		Type:  globals.FileType,
		Id:    req.FileId,
	})
	if err != nil {
		return err
	}

	return nil
}

// Search File handler
func (s *Service) Search(ctx context.Context, req *proto_file.SearchRequest, rsp *proto_file.SearchResponse) error {
	if err := validate.Exists(ctx, req.Index); err != nil {
		return err
	}

	res, err := operations.Search(ctx, &proto_operations.SearchRequest{
		Index:                req.Index,
		Term:                 req.Term,
		From:                 req.From,
		Size:                 req.Size,
		Category:             req.Category,
		Url:                  req.Url,
		Depth:                req.Depth,
		Type:                 globals.FileType,
		FileType:             req.FileType,
		Access:               req.Access,
		ContentCategory:      req.ContentCategory,
		NoKazoupFileOriginal: req.NoKazoupFileOriginal,
	})
	if err != nil {
		return errors.InternalServerError(globals.FILE_SERVICE_NAME, err.Error())
	}

	rsp.Result = res.Result
	rsp.Info = res.Info

	return nil
}

// Share file handler
func (s *Service) Share(ctx context.Context, req *proto_file.ShareRequest, rsp *proto_file.ShareResponse) error {
	if err := validate.Exists(ctx, req.OriginalId, req.DatasourceId); err != nil {
		return err
	}

	// Instantiate file system
	fsys, err := NewFileSystem(ctx, req.DatasourceId)
	if err != nil {
		return err
	}

	// Authorize datasource / refresh token
	auth, err := fsys.Authorize()
	if err != nil {
		return err
	}

	// Update token in DB
	if err := fs.UpdateFsAuth(ctx, req.DatasourceId, auth); err != nil {
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

	if _, err := operations.Update(ctx, &proto_operations.UpdateRequest{
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
