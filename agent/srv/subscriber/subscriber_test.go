package subscriber

import (
    "testing"
	"golang.org/x/net/context"
    save "github.com/kazoup/platform/lib/protomsg/saveremotefiles"
    "github.com/micro/go-micro"
)

var (
    agentServiceTaskHandler = new(AgentServiceTaskHandler)
)

// Announce subscriber LogAgentSaveMessage method unit test
func TestLogAgentSaveMessage(t *testing.T) {
    var LogAgentSaveMessageTestData = []struct {
        ctx context.Context
        msg *save.SaveRemoteFilesMessage
    } {
        {
            // Success
            micro.NewContext(ctx, srv),
            &save.SaveRemoteFilesMessage{
                Id: "a716f1408cfff7afd943acd45dcfa0a4",
                OriginalId: "id:lXWZMx78s2AAAAAAAAAAPw",
                OriginalDownloadRef: "",
                PreviewUrl: "/gopher.jpg",
                UserId: "google-apps|pablo.aguirre@kazoup.com",
                Name: "gopher.jpg",
                Url: "https://www.dropbox.com/home?preview=gopher.jpg",
                Modified: "2017-03-23T10:43:34Z",
                FileSize: 75679,
                IsDir: false,
                Category: "Pictures",
                MimeType: "",
                Depth: 0,
                FileType: "dropbox",
                LastSeen: 1496909821,
                Access: "private",
                DatasourceId: "e80d54ad29d18cb62cf9bb2bb54fcfd5",
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
