package elastic

import (
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/tjarratt/babble"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
)

func Test_elastic_Init(t *testing.T) {
	type fields struct {
		Client *elib.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			if err := e.Init(); (err != nil) != tt.wantErr {
				t.Errorf("elastic.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_elastic_Create(t *testing.T) {
	type fields struct {
		Client *elib.Client
	}
	type args struct {
		ctx context.Context
		req *proto_operations.CreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.CreateResponse
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("elastic.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_elastic_Read(t *testing.T) {
	type fields struct {
		Client *elib.Client
	}
	type args struct {
		ctx context.Context
		req *proto_operations.ReadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.ReadResponse
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Read(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("elastic.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_elastic_Update(t *testing.T) {
	type fields struct {
		Client *elib.Client
	}
	type args struct {
		ctx context.Context
		req *proto_operations.UpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.UpdateResponse
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("elastic.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_elastic_Delete(t *testing.T) {
	type fields struct {
		Client *elib.Client
	}
	type args struct {
		ctx context.Context
		req *proto_operations.DeleteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.DeleteResponse
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("elastic.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_elastic_DeleteByQuery(t *testing.T) {
	type fields struct {
		Client *elib.Client
	}
	type args struct {
		ctx context.Context
		req *proto_operations.DeleteByQueryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.DeleteByQueryResponse
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.DeleteByQuery(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.DeleteByQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("elastic.DeleteByQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_elastic_Search(t *testing.T) {
	type fields struct {
		Client *elib.Client
	}
	type args struct {
		ctx context.Context
		req *proto_operations.SearchRequest
	}
	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	// Client
	ec, err := elib.NewSimpleClient(
		elib.SetURL(url),
		elib.SetBasicAuth(username, password),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		t.Error(err)
	}
	//Random words generator
	babbler := babble.NewBabbler()
	//set your own separator
	babbler.Separator = " "
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.SearchResponse
		wantErr bool
	}{
		{"test-random-request",
			*&fields{ec},
			*&args{context.WithValue(
				context.TODO(),
				kazoup_context.UserIdCtxKey{},
				kazoup_context.UserIdCtxValue("google-apps|johan.holder@kazoup.com"),
			),
				&proto_operations.SearchRequest{
					Index:                "*",                   //"156cd7aee6bc4688cef74f02ce13bcee",
					Term:                 "sales and marketing", //babbler.Babble(),                   //String(5),
					From:                 0,
					Size:                 45,
					Type:                 "file",
					FileType:             "files",
					Depth:                0,
					NoKazoupFileOriginal: false,
				}},
			&proto_operations.SearchResponse{},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			t.Log(tt.args.req)

			got, err := e.Search(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got.GetInfo())

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("elastic.Search() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func BenchmarkElasticsearchSearch(b *testing.B) {
	// type fields struct {
	// 	Client *elib.Client
	// }
	// type args struct {
	// 	ctx context.Context
	// 	req *proto_operations.SearchRequest
	// }
	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	// Client
	ec, err := elib.NewSimpleClient(
		elib.SetURL(url),
		elib.SetBasicAuth(username, password),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		b.Error(err)
	}
	//Random words generator
	babbler := babble.NewBabbler()
	//set your own separator
	babbler.Separator = " "
	// type test struct {
	// 	name    string
	// 	fields  fields
	// 	args    args
	// 	want    *proto_operations.SearchResponse
	// 	wantErr bool
	// }
	e := &elastic{
		Client: ec,
	}
	ctx := context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue("google-apps|johan.holder@kazoup.com"),
	)
	for n := 0; n < b.N; n++ {

		req := &proto_operations.SearchRequest{
			Index:                "*",       //"156cd7aee6bc4688cef74f02ce13bcee",
			Term:                 String(5), //"kazoup", //babbler.Babble(), //String(5),
			From:                 0,
			Size:                 45,
			Type:                 "file",
			FileType:             "files",
			Depth:                0,
			NoKazoupFileOriginal: false,
		}
		//b.Log(req)
		_, err := e.Search(
			ctx,
			req,
		)
		if err != nil {
			b.Error(err)
		}
		//b.Log(res)
	}
}

const charset = "kazoup"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
