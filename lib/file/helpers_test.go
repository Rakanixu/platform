package file

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	gmailhelper "github.com/kazoup/platform/lib/gmail"
	"github.com/kazoup/platform/lib/onedrive"
	"github.com/kazoup/platform/lib/slack"
	googledrive "google.golang.org/api/drive/v3"
	"reflect"
	"testing"
)

var (
	datasourceId = "datasource_id"
	userId       = "user_id"
	index        = "index"
)

func TestNewFileFromString(t *testing.T) {
	slackFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.Slack,
	})
	if err != nil {
		t.Fatal(err)
	}

	googleDriveFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.GoogleDrive,
	})
	if err != nil {
		t.Fatal(err)
	}

	gmailFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.Gmail,
	})
	if err != nil {
		t.Fatal(err)
	}

	oneDriveFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.OneDrive,
	})
	if err != nil {
		t.Fatal(err)
	}

	dropboxFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.Dropbox,
	})
	if err != nil {
		t.Fatal(err)
	}

	boxFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.Box,
	})
	if err != nil {
		t.Fatal(err)
	}

	mockFileBytes, err := json.Marshal(KazoupFile{
		FileType: globals.Mock,
	})
	if err != nil {
		t.Fatal(err)
	}

	noFileBytes, err := json.Marshal(KazoupFile{
		FileType: "invalid_type",
	})
	if err != nil {
		t.Fatal(err)
	}

	type out struct {
		file File
		err  error
	}

	testData := []struct {
		in  string
		out out
	}{
		{
			string(slackFileBytes),
			out{
				&KazoupSlackFile{
					KazoupFile{
						FileType: globals.Slack,
					},
					nil,
				},
				nil,
			},
		},
		{
			string(googleDriveFileBytes),
			out{
				&KazoupGoogleFile{
					KazoupFile{
						FileType: globals.GoogleDrive,
					},
					nil,
				},
				nil,
			},
		},
		{
			string(gmailFileBytes),
			out{
				&KazoupGmailFile{
					KazoupFile{
						FileType: globals.Gmail,
					},
					nil,
				},
				nil,
			},
		},
		{
			string(oneDriveFileBytes),
			out{
				&KazoupOneDriveFile{
					KazoupFile{
						FileType: globals.OneDrive,
					},
					nil,
				},
				nil,
			},
		},
		{
			string(dropboxFileBytes),
			out{
				&KazoupDropboxFile{
					KazoupFile{
						FileType: globals.Dropbox,
					},
					nil,
				},
				nil,
			},
		},
		{
			string(boxFileBytes),
			out{
				&KazoupBoxFile{
					KazoupFile{
						FileType: globals.Box,
					},
					nil,
				},
				nil,
			},
		},
		{
			string(mockFileBytes),
			out{
				&KazoupMockFile{
					KazoupFile{
						FileType: globals.Mock,
					},
					nil,
				},
				nil,
			},
		},
		{
			string(noFileBytes),
			out{
				nil,
				errors.ErrInvalidFile,
			},
		},
	}

	for _, tt := range testData {
		result, err := NewFileFromString(tt.in)

		if !reflect.DeepEqual(tt.out.file, result) {
			t.Errorf("Expected: %v, got: %v", tt.out.file, result)
		}

		if tt.out.err != err {
			t.Errorf("Expected error: %v, got: %v", tt.out.err, err)
		}
	}
}

func TestNewKazoupFileFromGoogleDriveFile(t *testing.T) {
	type in struct {
		original     googledrive.File
		datasourceId string
		userId       string
		index        string
	}

	testData := []struct {
		in  in
		out *KazoupGoogleFile
	}{
		{
			in{
				googledrive.File{
					Id: "googledrive_id",
				},
				datasourceId,
				userId,
				index,
			},
			&KazoupGoogleFile{
				KazoupFile: KazoupFile{
					DatasourceId: datasourceId,
					UserId:       userId,
					Index:        index,
					Access:       globals.ACCESS_PRIVATE,
				},
				Original: &googledrive.File{
					Id: "googledrive_id",
				},
			},
		},
		{
			in{
				googledrive.File{
					Id:      "googledrive_id",
					Trashed: true,
				},
				datasourceId,
				userId,
				index,
			},
			nil,
		},
	}

	for _, tt := range testData {
		result := NewKazoupFileFromGoogleDriveFile(
			tt.in.original,
			tt.in.datasourceId,
			tt.in.userId,
			tt.in.index,
		)

		if reflect.TypeOf(result) != reflect.TypeOf(tt.out) {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}

func TestNewKazoupFileFromSlackFile(t *testing.T) {
	type in struct {
		original     slack.SlackFile
		datasourceId string
		userId       string
		index        string
	}

	testData := []struct {
		in  in
		out *KazoupSlackFile
	}{
		{
			in{
				slack.SlackFile{
					ID: "slack_id",
				},
				datasourceId,
				userId,
				index,
			},
			&KazoupSlackFile{
				KazoupFile: KazoupFile{
					DatasourceId: datasourceId,
					UserId:       userId,
					Index:        index,
					Access:       globals.ACCESS_PRIVATE,
				},
				Original: &slack.SlackFile{
					ID: "slack_id",
				},
			},
		},
	}

	for _, tt := range testData {
		result := NewKazoupFileFromSlackFile(
			tt.in.original,
			tt.in.datasourceId,
			tt.in.userId,
			tt.in.index,
		)

		if reflect.TypeOf(result) != reflect.TypeOf(tt.out) {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}

func TestNewKazoupFileFromOneDriveFile(t *testing.T) {
	type in struct {
		original     onedrive.OneDriveFile
		datasourceId string
		userId       string
		index        string
	}

	testData := []struct {
		in  in
		out *KazoupOneDriveFile
	}{
		{
			in{
				onedrive.OneDriveFile{
					ID: "onedrive_id",
				},
				datasourceId,
				userId,
				index,
			},
			&KazoupOneDriveFile{
				KazoupFile: KazoupFile{
					DatasourceId: datasourceId,
					UserId:       userId,
					Index:        index,
					Access:       globals.ACCESS_PRIVATE,
				},
				Original: &onedrive.OneDriveFile{
					ID: "onedrive_id",
				},
			},
		},
	}

	for _, tt := range testData {
		result := NewKazoupFileFromOneDriveFile(
			tt.in.original,
			tt.in.datasourceId,
			tt.in.userId,
			tt.in.index,
		)

		if reflect.TypeOf(result) != reflect.TypeOf(tt.out) {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}

func TestNewKazoupFileFromDropboxFile(t *testing.T) {
	type in struct {
		original     dropbox.DropboxFile
		datasourceId string
		userId       string
		index        string
	}

	testData := []struct {
		in  in
		out *KazoupDropboxFile
	}{
		{
			in{
				dropbox.DropboxFile{
					ID: "dropbox_id",
				},
				datasourceId,
				userId,
				index,
			},
			&KazoupDropboxFile{
				KazoupFile: KazoupFile{
					DatasourceId: datasourceId,
					UserId:       userId,
					Index:        index,
					Access:       globals.ACCESS_PRIVATE,
				},
				Original: &dropbox.DropboxFile{
					ID: "dropbox_id",
				},
			},
		},
		{
			in{
				dropbox.DropboxFile{
					ID:         "dropbox_id",
					DropboxTag: "deleted",
				},
				datasourceId,
				userId,
				index,
			},
			nil,
		},
	}

	for _, tt := range testData {
		result := NewKazoupFileFromDropboxFile(
			tt.in.original,
			tt.in.datasourceId,
			tt.in.userId,
			tt.in.index,
		)

		if reflect.TypeOf(result) != reflect.TypeOf(tt.out) {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}

func TestNewKazoupFileFromBoxFile(t *testing.T) {
	type in struct {
		original     box.BoxFileMeta
		datasourceId string
		userId       string
		index        string
	}

	testData := []struct {
		in  in
		out *KazoupBoxFile
	}{
		{
			in{
				box.BoxFileMeta{
					ID:         "box_id",
					ModifiedAt: "2006-01-02T15:04:05Z",
				},
				datasourceId,
				userId,
				index,
			},
			&KazoupBoxFile{
				KazoupFile: KazoupFile{
					DatasourceId: datasourceId,
					UserId:       userId,
					Index:        index,
					Access:       globals.ACCESS_PRIVATE,
				},
				Original: &box.BoxFileMeta{
					ID: "box_id",
				},
			},
		},
	}

	for _, tt := range testData {
		result := NewKazoupFileFromBoxFile(
			tt.in.original,
			tt.in.datasourceId,
			tt.in.userId,
			tt.in.index,
		)

		if reflect.TypeOf(result) != reflect.TypeOf(tt.out) {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}

func TestNewKazoupFileFromGmailFile(t *testing.T) {
	type in struct {
		original     gmailhelper.GmailFile
		datasourceId string
		userId       string
		index        string
	}

	testData := []struct {
		in  in
		out *KazoupGmailFile
	}{
		{
			in{
				gmailhelper.GmailFile{
					Id: "gmail_id",
				},
				datasourceId,
				userId,
				index,
			},
			&KazoupGmailFile{
				KazoupFile: KazoupFile{
					DatasourceId: datasourceId,
					UserId:       userId,
					Index:        index,
					Access:       globals.ACCESS_PRIVATE,
				},
				Original: &gmailhelper.GmailFile{
					Id: "gmail_id",
				},
			},
		},
	}

	for _, tt := range testData {
		result := NewKazoupFileFromGmailFile(
			tt.in.original,
			tt.in.datasourceId,
			tt.in.userId,
			"dsURL",
			tt.in.index,
		)

		if reflect.TypeOf(result) != reflect.TypeOf(tt.out) {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}

func TestNewKazoupFileFromMockFile(t *testing.T) {
	testData := []struct {
		out *KazoupMockFile
	}{
		{
			&KazoupMockFile{
				KazoupFile: KazoupFile{
					DatasourceId: datasourceId,
					UserId:       userId,
					Index:        index,
					Access:       globals.ACCESS_PRIVATE,
				},
			},
		},
	}

	for _, tt := range testData {
		result := NewKazoupFileFromMockFile()

		if reflect.TypeOf(result) != reflect.TypeOf(tt.out) {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}
