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

func TestOnDatasourceCreate(t *testing.T) {
	e := &proto_datasource.CreateRequest{
		&proto_datasource.Endpoint{
			UserId: TEST_USER_ID,
		},
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onDatasourceCreateTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_DATASOURCE_CREATE,
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
				Handler: globals.HANDLER_DATASOURCE_CREATE,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onDatasourceCreateTestData {
		err := announceHandler.OnDatasourceCreate(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnDatasourceDelete(t *testing.T) {
	e := &proto_datasource.DeleteRequest{
		Id:    "test",
		Index: "test",
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onDatasourceDeleteTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_DATASOURCE_DELETE,
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
				Handler: globals.HANDLER_DATASOURCE_DELETE,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onDatasourceDeleteTestData {
		err := announceHandler.OnDatasourceDelete(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnDatasourceScan(t *testing.T) {
	e := &proto_datasource.ScanRequest{
		Id: "test",
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onDatasourceScanTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_DATASOURCE_SCAN,
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
				Handler: globals.HANDLER_DATASOURCE_SCAN,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onDatasourceScanTestData {
		err := announceHandler.OnDatasourceScan(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnDatasourceScanAll(t *testing.T) {
	withId := &proto_datasource.ScanAllRequest{
		DatasourcesId: []string{"test"},
	}
	withoutId := &proto_datasource.ScanAllRequest{
		DatasourcesId: []string{},
	}

	b1, err := json.Marshal(withId)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := json.Marshal(withoutId)
	if err != nil {
		t.Fatal(err)
	}

	var onDatasourceScanAllTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_DATASOURCE_SCANALL,
				Data:    string(b1),
			},
			nil,
		},
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_DATASOURCE_SCANALL,
				Data:    string(b2),
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
				Handler: globals.HANDLER_DATASOURCE_SCANALL,
				Data:    string(b1),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onDatasourceScanAllTestData {
		err := announceHandler.OnDatasourceScanAll(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}
