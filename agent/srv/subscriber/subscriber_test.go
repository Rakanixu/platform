package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/file"
	save "github.com/kazoup/platform/lib/protomsg/saveremotefiles"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"testing"
	"time"
)

var (
	agentServiceTaskHandler = new(AgentServiceTaskHandler)
)

// Announce subscriber LogAgentSaveMessage method unit test
func TestLogAgentSaveMessage(t *testing.T) {
	// Parse string to Time
	m, err := time.Parse(time.RFC3339, "2017-03-23T10:43:34Z")
	if err != nil {
		t.Fatal(err)
	}

	// Testing KazoupFile
	f := file.KazoupFile{
		ID:                  "a716f1408cfff7afd943acd45dcfa0a4",
		OriginalID:          "id:lXWZMx78s2AAAAAAAAAAPw",
		OriginalDownloadRef: "",
		PreviewUrl:          "/gopher.jpg",
		UserId:              "google-apps|pablo.aguirre@kazoup.com",
		Name:                "gopher.jpg",
		URL:                 "https://www.dropbox.com/home?preview=gopher.jpg",
		Modified:            m,
		FileSize:            75679,
		IsDir:               false,
		Category:            "Pictures",
		MimeType:            "",
		Depth:               0,
		FileType:            "dropbox",
		LastSeen:            1496909821,
		Access:              "private",
		DatasourceId:        "e80d54ad29d18cb62cf9bb2bb54fcfd5",
		Index:               "index8d68d0671dfb-4201-bb8b-ee3dba1cc3ff",
	}

	// Marshal file structure to JSON
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	var LogAgentSaveMessageTestData = []struct {
		ctx context.Context
		msg *save.SaveRemoteFilesMessage
	}{
		{
			// Success
			micro.NewContext(ctx, srv),
			&save.SaveRemoteFilesMessage{
				UserId: f.UserId,
				Index:  f.Index,
				Data:   string(b),
			},
		},
	}

	for _, tt := range LogAgentSaveMessageTestData {
		err := agentServiceTaskHandler.LogAgentSaveMessage(tt.ctx, tt.msg)
		if err != nil {
			t.Fatal(err)
		}
	}
}
