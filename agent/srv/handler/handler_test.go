package handler

import (
    "github.com/kazoup/platform/agent/srv/proto/agent"
	"golang.org/x/net/context"
	kazoup_context "github.com/kazoup/platform/lib/context"
    "testing"
	"github.com/micro/go-micro/metadata"
    _ "github.com/kazoup/platform/lib/quota/mock"
)

// Constants
const (
    TEST_USER_ID = "test_user"
    TEST_JSON_OBJECT = `
		{
			"id" : "a716f1408cfff7afd943acd45dcfa0a4",
            "original_id" : "id:lXWZMx78s2AAAAAAAAAAPw",
			"original_download_ref" : "",
			"preview_url" : "/gopher.jpg",
			"user_id": "google-apps|pablo.aguirre@kazoup.com",
			"name" : "gopher.jpg",
            "url" : "https://www.dropbox.com/home?preview=gopher.jpg",
            "modified": "2017-03-23T10:43:34Z",
			"file_size" : 75679,
			"is_dir" : false,
			"category" : "Pictures",
			"mime_type" : "",
			"depth" : 0,
			"file_type" : "dropbox",
			"last_seen" : 1496909821,
			"access" : "private",
			"datasource_id" : "e80d54ad29d18cb62cf9bb2bb54fcfd5",
			"index" : "index8d68d0671dfb-4201-bb8b-ee3dba1cc3ff"
		}
    `
)

// Helping variables
var (
    srv = new(Service)
    ctx = context.WithValue(
        context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
    )
)

// Handler Save method unit test
func TestSave(t *testing.T) {
    // Test data
    var saveTestData = []struct{
		ctx         context.Context
		req         *proto_agent.SaveRequest
		expectedRsp *proto_agent.SaveResponse
		rsp         *proto_agent.SaveResponse
    }{
        {
            // Quota has been exceeded
            metadata.NewContext(ctx, map[string]string{
                "Quota-Exceeded": "true",
            }),
            &proto_agent.SaveRequest{
                Data: TEST_JSON_OBJECT,
            },
            &proto_agent.SaveResponse{
                Info: QUOTA_EXCEEDED_MSG,
            },
            &proto_agent.SaveResponse{},
        },
        {
            // Quota has been exceeded
            metadata.NewContext(ctx, map[string]string{
                "Quota-Exceeded": "false",
            }),
            &proto_agent.SaveRequest{
                Data: TEST_JSON_OBJECT,
            },
            &proto_agent.SaveResponse{
                Info: "",
            },
            &proto_agent.SaveResponse{},
        },
    }

    // Run tests and check responses
    for _, tt := range saveTestData {
        if err := srv.Save(tt.ctx, tt.req, tt.rsp); err != nil {
            t.Fatal(err)
        }

        if tt.expectedRsp.Info != tt.rsp.Info {
            t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Info, tt.rsp.Info)
        }
    }
}
