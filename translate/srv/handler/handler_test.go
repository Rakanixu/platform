package handler

import (
	"reflect"
	"testing"

	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/translate/srv/proto/translate"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
)

const (
	testUserID = "test_user"
)

var (
	srv = new(Service)
	ctx = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(testUserID),
	)
)

func TestTranslate(t *testing.T) {
	var translateTestData = []struct {
		ctx         context.Context
		req         *proto_translate.TranslateRequest
		expectedRsp *proto_translate.TranslateResponse
		rsp         *proto_translate.TranslateResponse
	}{
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "true",
			}),
			&proto_translate.TranslateRequest{
				Text: []string{
					"Hello world!",
				},
				SourceLang: "en",
				DestLang:   "fr",
			},
			&proto_translate.TranslateResponse{
				Info: quotaExceededMsg,
			},
			&proto_translate.TranslateResponse{},
		},
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "false",
			}),
			&proto_translate.TranslateRequest{
				Text: []string{
					"Hello world!",
				},
				SourceLang: "en",
				DestLang:   "fr",
			},
			&proto_translate.TranslateResponse{
				Info: "",
				Translations: []string{
					"Bonjour le monde!",
				},
			},
			&proto_translate.TranslateResponse{},
		},
	}

	for _, tt := range translateTestData {
		if err := srv.Translate(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tt.rsp, tt.expectedRsp) {
			t.Error("Expected:", tt.expectedRsp, "got:", tt.rsp)
		}
	}
}
