package cloudstorage

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"reflect"
	"testing"
)

func TestNewCloudStorageFromEndpoint(t *testing.T) {
	type out struct {
		result CloudStorage
		err    error
	}

	var testData = []struct {
		endpoint  *proto_datasource.Endpoint
		connector string
		out       out
	}{
		{
			&proto_datasource.Endpoint{},
			globals.Slack,
			out{
				NewSlackCloudStorage(&proto_datasource.Endpoint{}),
				nil,
			},
		},
		{
			&proto_datasource.Endpoint{},
			globals.GoogleDrive,
			out{
				NewGoogleDriveCloudStorage(&proto_datasource.Endpoint{}),
				nil,
			},
		},
		{
			&proto_datasource.Endpoint{},
			globals.Gmail,
			out{
				NewGmailCloudStorage(&proto_datasource.Endpoint{}),
				nil,
			},
		},
		{
			&proto_datasource.Endpoint{},
			globals.OneDrive,
			out{
				NewOneDriveCloudStorage(&proto_datasource.Endpoint{}),
				nil,
			},
		},
		{
			&proto_datasource.Endpoint{},
			globals.Dropbox,
			out{
				NewDropboxCloudStorage(&proto_datasource.Endpoint{}),
				nil,
			},
		},
		{
			&proto_datasource.Endpoint{},
			globals.Box,
			out{
				NewBoxCloudStorage(&proto_datasource.Endpoint{}),
				nil,
			},
		},
	}

	for _, tt := range testData {
		result, err := NewCloudStorageFromEndpoint(tt.endpoint, tt.connector)

		if !reflect.DeepEqual(tt.out.result, result) {
			t.Errorf("Expected: %v, got: %v", tt.out.result, result)
		}

		if tt.out.err != err {
			t.Errorf("Expected error: %v, got error: %v", tt.out.err, err)
		}
	}
}
