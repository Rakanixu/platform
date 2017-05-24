package main

/*
import (
	"github.com/kazoup/platform/audio/srv/handler"
	"github.com/kazoup/platform/audio/srv/proto/audio"
	_ "github.com/kazoup/platform/lib/quota/mock"
	"github.com/kazoup/platform/lib/wrappers"
	micro "github.com/micro/go-micro"
	b "github.com/micro/go-micro/broker/mock"
	"github.com/micro/go-micro/metadata"
	r "github.com/micro/go-micro/registry/mock"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"testing"
)

const (
	SRV_NAME = "com.kazoup.srv.audio"
)

var (
	service micro.Service
	client  proto_audio.ServiceClient
)

*/
/*func subscriber(ctx context.Context, msg *example.Message) error {

	var err error
	if msg.GetSay() != "Hi" {
		log.Print("error")
		err = fmt.Errorf("Expextected Hi got %s", msg.GetSay())
		os.Exit(1)
	}
	close(waitForMsg)
	return err
}*/ /*


func TestMain(m *testing.M) {
	// Tests Setup

	//wait channels
	waitForService := make(chan bool)
	//make service
	service = micro.NewService(
		micro.Name(SRV_NAME),
		micro.Registry(r.NewRegistry()),
		micro.Broker(b.NewBroker()),
		micro.AfterStart(func() error {
			close(waitForService)
			return nil
		}),
	)

	service.Init(
		//micro.WrapClient(wrappers.ContextClientWrapper(service)),
		micro.WrapHandler(
			wrappers.UserTestHandlerWrapper,
			wrappers.NewContextHandlerWrapper(service),
			wrappers.NewAfterHandlerWrapper(),
		),
		micro.WrapSubscriber(
			wrappers.UserTestSubscriberWrapper,
			wrappers.NewContextSubscriberWrapper(service),
			wrappers.NewAfterSubscriberWrapper(),
		),
	)

	// Register handler to be tested
	proto_audio.RegisterServiceHandler(service.Server(), new(handler.Service))

	client = proto_audio.NewServiceClient(SRV_NAME, service.Client())

	// Register function as subscriber
*/
/*	err := service.Server().Subscribe(
		service.Server().NewSubscriber("topic.go.micro.srv.test", subscriber),
	)
	if err != nil {
		os.Exit(1)
	}*/ /*

	//start service TODO handle returned error
	go service.Run()

	//wait for start
	<-waitForService

	exitVal := m.Run()

	// Tear down
	if err := service.Server().Deregister(); err != nil {
		os.Exit(1)
	}

	if err := service.Server().Stop(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitVal)
}

func TestHandlerEnrichFile(t *testing.T) {
	var enrinchFilesTestData = []struct {
		ctx context.Context
		req *proto_audio.EnrichFileRequest
		rsp *proto_audio.EnrichFileResponse
	}{
		// Quota has been excedded
		{
			metadata.NewContext(context.TODO(), map[string]string{
				"Quota-Exceeded": "true",
			}),
			&proto_audio.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_audio.EnrichFileResponse{
				Info: handler.QUOTA_EXCEEDED_MSG,
			},
		},
		// Quota has not been exceeded
		{
			metadata.NewContext(context.TODO(), map[string]string{
				"Quota-Exceeded": "false",
			}),
			&proto_audio.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_audio.EnrichFileResponse{
				Info: "",
			},
		},
	}

	for _, tt := range enrinchFilesTestData {
		rsp, err := client.EnrichFile(tt.ctx, tt.req)
		if err != nil {
			t.Fatal(err)
		}

		if rsp.Info != tt.rsp.Info {
			t.Errorf("Expected '%v', got: '%v'", tt.rsp.Info, rsp.Info)
		}
	}
}

func TestHandlerEnrichDatasource(t *testing.T) {
	var enrinchDatasourceTestData = []struct {
		ctx context.Context
		req *proto_audio.EnrichDatasourceRequest
		rsp *proto_audio.EnrichDatasourceResponse
	}{
		{
			context.TODO(),
			&proto_audio.EnrichDatasourceRequest{
				Id: "test_id",
			},
			&proto_audio.EnrichDatasourceResponse{},
		},
	}

	for _, tt := range enrinchDatasourceTestData {
		rsp, err := client.EnrichDatasource(tt.ctx, tt.req)
		if err != nil {
			t.Fatal(err)
		}

		if rsp.Info != tt.rsp.Info {
			t.Errorf("Expected '%v', got: '%v'", tt.rsp.Info, rsp.Info)
		}
	}
}

func TestHandlerHealth(t *testing.T) {
	rsp, err := client.Health(context.TODO(), &proto_audio.HealthRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if rsp.Status != http.StatusOK {
		t.Errorf("Expected '%v', got: '%v'", http.StatusOK, rsp.Status)
	}
}

//TestTestServiceSubscriber tests subscribers
*/
/*func TestTestServiceSubscriber(t *testing.T) {

	// //Test publish subscribe
	publisher := example.NewPublisher("topic.go.micro.srv.test", service.Client())

	err := publisher.Publish(context.TODO(), &example.Message{
		Say: "Hi",
	})
	if err != nil {
		t.Error(err)
	}
	//Wait for msg
	<-waitForMsg

}*/
