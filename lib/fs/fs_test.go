package fs

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

func TestNewFsFromEndpoint(t *testing.T) {
	type out struct {
		typ interface{}
		err error
	}

	testData := []struct {
		in  *proto_datasource.Endpoint
		out out
	}{
		{
			&proto_datasource.Endpoint{
				Url: globals.Slack,
			},
			out{
				typ: "*fs.SlackFs",
				err: nil,
			},
		},
		{
			&proto_datasource.Endpoint{
				Url: globals.GoogleDrive,
			},
			out{
				typ: "*fs.GoogleDriveFs",
				err: nil,
			},
		},
		{
			&proto_datasource.Endpoint{
				Url: globals.Gmail,
			},
			out{
				typ: "*fs.GmailFs",
				err: nil,
			},
		},
		{
			&proto_datasource.Endpoint{
				Url: globals.OneDrive,
			},
			out{
				typ: "*fs.OneDriveFs",
				err: nil,
			},
		},
		{
			&proto_datasource.Endpoint{
				Url: globals.Dropbox,
			},
			out{
				typ: "*fs.DropboxFs",
				err: nil,
			},
		},
		{
			&proto_datasource.Endpoint{
				Url: globals.Box,
			},
			out{
				typ: "*fs.BoxFs",
				err: nil,
			},
		},
		{
			&proto_datasource.Endpoint{
				Url: globals.Mock,
			},
			out{
				typ: "*fs.MockFs",
				err: nil,
			},
		},
		{
			&proto_datasource.Endpoint{
				Url: "invalid",
			},
			out{
				typ: nil,
				err: errors.ErrInvalidFileSystem,
			},
		},
	}

	for _, tt := range testData {
		result, err := NewFsFromEndpoint(tt.in)

		if (result != nil) && tt.out.typ != reflect.TypeOf(result).String() {
			t.Errorf("Expected %v, got %v", tt.out.typ, reflect.TypeOf(result).String())
		}

		if tt.out.err != err {
			t.Errorf("Expected %v, got %v", tt.out.err, err)
		}
	}
}

func TestUpdateFsAuth(t *testing.T) {
	ctx := context.TODO()
	id := "id"
	token := &proto_datasource.Token{}

	result := UpdateFsAuth(ctx, id, token)

	if nil != result {
		t.Errorf("Expected %v, got %v", nil, result)
	}
}

func TestClearIndex(t *testing.T) {
	ctx := context.TODO()
	endpoint := &proto_datasource.Endpoint{
		Url: globals.Mock,
	}

	result := ClearIndex(ctx, endpoint)

	if nil != result {
		t.Errorf("Expected %v, got %v", nil, result)
	}
}
