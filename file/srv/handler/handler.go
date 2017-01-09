package handler

import (
	"encoding/json"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
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

	// Get userId for later
	uId, err := globals.ParseJWTToken(ctx)
	if err != nil {
		return err
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

	// Update token in DB
	if err := UpdateFileSystemAuth(f.Client, ctx, req.DatasourceId, auth); err != nil {
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
	_, err = DeleteFromDB(f.Client, ctx, &db_proto.DeleteRequest{
		Index: req.Index,
		Type:  globals.FileType,
		Id:    req.FileId,
	})
	if err != nil {
		return err
	}

	// Publish notification topic, let client know when to refresh itself
	if err := f.Client.Publish(globals.NewSystemContext(), f.Client.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
		Method: globals.NOTIFY_REFRESH_SEARCH,
		UserId: uId,
	})); err != nil {
		log.Print("Publishing (notify file) error %s", err)
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
