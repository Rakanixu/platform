package elastic

import (
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"encoding/json"
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/tjarratt/babble"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
)

const (
	USER_ID          = "google-apps|johan.holder@kazoup.com"
	TYPE_FILE        = "file"
	INDEX_TEST_FILES = "index_elastic_test"
)

var (
	file_id   string
	test_file *file.KazoupFile
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
	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	//Random words generator
	babbler := babble.NewBabbler()

	// Client
	ec, err := elib.NewSimpleClient(
		elib.SetURL(url),
		elib.SetBasicAuth(username, password),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		t.Fatal(err)
	}

	test_file = &file.KazoupFile{
		Name:     babbler.Babble(),
		ID:       babbler.Babble(),
		Content:  babbler.Babble(),
		LastSeen: time.Now().Unix(),
	}
	file_id = test_file.ID

	b, err := json.Marshal(test_file)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.CreateResponse
		wantErr bool
	}{
		{
			"Create",
			fields{ec},
			args{
				context.WithValue(
					context.TODO(),
					kazoup_context.UserIdCtxKey{},
					kazoup_context.UserIdCtxValue(USER_ID),
				),
				&proto_operations.CreateRequest{
					Index: INDEX_TEST_FILES,
					Type:  TYPE_FILE,
					Id:    test_file.ID,
					Data:  string(b),
				},
			},
			&proto_operations.CreateResponse{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("elastic.Create() error = %v, wantErr %v", err, tt.wantErr)
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
		t.Fatal(err)
	}

	b, err := json.Marshal(test_file)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.ReadResponse
		wantErr bool
	}{
		{
			"Read",
			fields{ec},
			args{
				context.WithValue(
					context.TODO(),
					kazoup_context.UserIdCtxKey{},
					kazoup_context.UserIdCtxValue(USER_ID),
				),
				&proto_operations.ReadRequest{
					Index: INDEX_TEST_FILES,
					Type:  TYPE_FILE,
					Id:    file_id,
				},
			},
			&proto_operations.ReadResponse{
				Result: string(b),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Read(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("elastic.Read() error = %v, wantErr %v", err, tt.wantErr)
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
		t.Fatal(err)
	}

	test_file.Modified = time.Now()
	test_file.FileType = globals.GoogleDrive

	b, err := json.Marshal(test_file)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.UpdateResponse
		wantErr bool
	}{
		{
			"Update",
			fields{ec},
			args{
				context.WithValue(
					context.TODO(),
					kazoup_context.UserIdCtxKey{},
					kazoup_context.UserIdCtxValue(USER_ID),
				),
				&proto_operations.UpdateRequest{
					Index: INDEX_TEST_FILES,
					Type:  TYPE_FILE,
					Id:    file_id,
					Data:  string(b),
				},
			},
			&proto_operations.UpdateResponse{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("elastic.Update() error = %v, wantErr %v", err, tt.wantErr)
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
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.DeleteResponse
		wantErr bool
	}{
		{
			"Delete",
			fields{ec},
			args{
				context.WithValue(
					context.TODO(),
					kazoup_context.UserIdCtxKey{},
					kazoup_context.UserIdCtxValue(USER_ID),
				),
				&proto_operations.DeleteRequest{
					Index: INDEX_TEST_FILES,
					Type:  TYPE_FILE,
					Id:    file_id,
				},
			},
			&proto_operations.DeleteResponse{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("elastic.Delete() error = %v, wantErr %v", err, tt.wantErr)
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
		t.Fatal(err)
	}

	// Create file
	Test_elastic_Create(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto_operations.DeleteByQueryResponse
		wantErr bool
	}{
		{
			"DeleteByQuery",
			fields{ec},
			args{
				context.WithValue(
					context.TODO(),
					kazoup_context.UserIdCtxKey{},
					kazoup_context.UserIdCtxValue(USER_ID),
				),
				&proto_operations.DeleteByQueryRequest{
					Indexes:  []string{INDEX_TEST_FILES},
					Types:    []string{TYPE_FILE},
					FileType: TYPE_FILE,
					LastSeen: test_file.LastSeen - 1,
				},
			},
			&proto_operations.DeleteByQueryResponse{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &elastic{
				Client: tt.fields.Client,
			}
			got, err := e.DeleteByQuery(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("elastic.DeleteByQuery() error = %v, wantErr %v", err, tt.wantErr)
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
		{
			"test-random-request",
			fields{ec},
			args{
				context.WithValue(
					context.TODO(),
					kazoup_context.UserIdCtxKey{},
					kazoup_context.UserIdCtxValue(USER_ID),
				),
				&proto_operations.SearchRequest{
					Index:                "*",                   //"156cd7aee6bc4688cef74f02ce13bcee",
					Term:                 "sales and marketing", //babbler.Babble(),                   //String(5),
					From:                 0,
					Size:                 45,
					Type:                 TYPE_FILE,
					FileType:             "files",
					Depth:                0,
					NoKazoupFileOriginal: false,
				},
			},
			&proto_operations.SearchResponse{},
			false,
		},
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
			}
			t.Log(got.GetInfo())
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
		kazoup_context.UserIdCtxValue(USER_ID),
	)
	for n := 0; n < b.N; n++ {
		req := &proto_operations.SearchRequest{
			Index:                "*",       //"156cd7aee6bc4688cef74f02ce13bcee",
			Term:                 String(5), //"kazoup", //babbler.Babble(), //String(5),
			From:                 0,
			Size:                 45,
			Type:                 TYPE_FILE,
			FileType:             "files",
			Depth:                0,
			NoKazoupFileOriginal: false,
		}

		_, err := e.Search(
			ctx,
			req,
		)
		if err != nil {
			b.Error(err)
		}
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
