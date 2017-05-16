package handler

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/objectstorage/mock"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type request struct{}

func (r request) Service() string {
	return "service"
}

func (r request) Method() string {
	return "method"
}

func (r request) ContentType() string {
	return "Content-Type"
}

func (r request) Request() interface{} {
	return r
}

func (r request) Stream() bool {
	return false
}

type response struct{}

type publication struct{}

func (p publication) Topic() string {
	return "topic"
}

func (p publication) Message() interface{} {
	type msg struct{}

	return new(msg)
}

func (p publication) ContentType() string {
	return "Content-Type"
}

type stream struct{}

func (s stream) Context() context.Context {
	return context.TODO()
}
func (s stream) Request() client.Request {
	return request{}
}

func (s stream) Send(interface{}) error {
	return nil
}

func (s stream) Recv(interface{}) error {
	return nil
}

func (s stream) Error() error {
	return nil
}

func (s stream) Close() error {
	return nil
}

type mockClient struct{}

func (c mockClient) Init(...client.Option) error {
	return nil
}

func (c mockClient) Options() client.Options {
	return client.Options{}
}

func (c mockClient) NewPublication(topic string, msg interface{}) client.Publication {
	return publication{}
}

func (c mockClient) NewRequest(service, method string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return request{}
}

func (c mockClient) NewProtoRequest(service, method string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return request{}
}

func (c mockClient) NewJsonRequest(service, method string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return request{}
}

func (c mockClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return nil
}

func (c mockClient) CallRemote(ctx context.Context, addr string, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return nil
}

func (c mockClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Streamer, error) {
	return stream{}, nil
}

func (c mockClient) StreamRemote(ctx context.Context, addr string, req client.Request, opts ...client.CallOption) (client.Streamer, error) {
	return stream{}, nil
}

func (c mockClient) Publish(ctx context.Context, p client.Publication, opts ...client.PublishOption) error {
	return nil
}

func (c mockClient) String() string {
	return ""
}

func TestSaveDatasource(t *testing.T) {
	result := SaveDatasource(context.TODO(), mockClient{}, "test_user", "url", &oauth2.Token{})

	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}

func TestSaveTmpToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test_user",
		"exp": time.Date(2050, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// RPC API ClientID can be found in https://manage.auth0.com/#/clients
	decoded, _ := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
	tokenString, err := token.SignedString(decoded)
	if err != nil {
		t.Fatalf("Error generating valid JWT %v", err)
	}

	result := SaveTmpToken("uuid", tokenString)

	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}

func TestRetrieveUserAndContextFromUUID(t *testing.T) {
	expectedUserId := "test_user"

	userId, _, err := RetrieveUserAndContextFromUUID("test_uuid")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedUserId != userId {
		t.Errorf("Expected %v, got %v", expectedUserId, userId)
	}
}

func TestPublishNotification(t *testing.T) {
	result := PublishNotification(context.TODO(), "test_user")

	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}

func TestCloseBrowserWindow(t *testing.T) {
	expected := `
		<script>
		'use stric';
			(function() {
				window.close();
			}());
		</script>
	`

	req, err := http.NewRequest(http.MethodGet, "/aaaa/asd", nil)
	if err != nil {
		t.Fatalf("Error building request: %v", err)
	}

	r := httptest.NewRecorder()
	handler := http.HandlerFunc(CloseBrowserWindow)
	handler.ServeHTTP(r, req)

	if expected != r.Body.String() {
		t.Errorf("Expected %v, got %v", expected, r.Body.String())
	}
}
