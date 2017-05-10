package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"testing"
)

const (
	TEST_USER_ID = "test_user"
)

var (
	announceHandler = new(AnnounceHandler)
	srv             = wrappers.NewKazoupService("test-service")
	ctx             = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

func TestOnDiscover(t *testing.T) {
	e := &proto_datasource.Endpoint{
		UserId: TEST_USER_ID,
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onDiscoverTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.DiscoverTopic,
				Data:    string(b),
			},
			nil,
		},
		// Ignore msg due to topic
		{
			ctx,
			&announce.AnnounceMessage{
				Handler: "ignore-me",
			},
			nil,
		},
		//Invalid context
		{
			ctx,
			&announce.AnnounceMessage{
				Handler: globals.DiscoverTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onDiscoverTestData {
		err := announceHandler.OnDiscover(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}
