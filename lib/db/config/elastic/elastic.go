package elastic

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/db/config"
	"github.com/kazoup/platform/lib/db/config/proto/config"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
	"os"
)

type elastic struct {
	Client *elib.Client
}

func init() {
	config.Register(new(elastic))
}

// Init config db
func (e *elastic) Init() error {
	var err error

	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://localhost:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	// Client
	e.Client, err = elib.NewSimpleClient(
		elib.SetURL(url),
		elib.SetBasicAuth(username, password),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		return err
	}

	return nil
}

// CreateIndex creates an ES index with template settings if match index* naming
func (e *elastic) CreateIndex(ctx context.Context, req *proto_config.CreateIndexRequest) (*proto_config.CreateIndexResponse, error) {
	exists, err := e.Client.IndexExists(req.Index).Do(ctx)
	if err != nil {
		return &proto_config.CreateIndexResponse{}, err
	}

	if !exists {
		// Create a new index.
		_, err := e.Client.CreateIndex(req.Index).Do(ctx)
		if err != nil {
			return &proto_config.CreateIndexResponse{}, err
		}
	}

	return &proto_config.CreateIndexResponse{}, nil
}

// Status elasticsearch cluster
func (e *elastic) Status(ctx context.Context, req *proto_config.StatusRequest) (*proto_config.StatusResponse, error) {
	cs, err := e.Client.ClusterState().Do(ctx)
	if err != nil {
		return &proto_config.StatusResponse{}, err
	}

	b, err := json.Marshal(cs)
	if err != nil {
		return &proto_config.StatusResponse{}, err
	}

	response := &proto_config.StatusResponse{
		Status: string(b),
	}

	return response, err
}

// AddAlias to assign indexes (aliases) per datasource
func (e *elastic) AddAlias(ctx context.Context, req *proto_config.AddAliasRequest) (*proto_config.AddAliasResponse, error) {
	_, err := e.Client.Alias().Add(req.Index, req.Alias).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &proto_config.AddAliasResponse{}, nil
}

// DeleteIndex from ES
func (e *elastic) DeleteIndex(ctx context.Context, req *proto_config.DeleteIndexRequest) (*proto_config.DeleteIndexResponse, error) {
	_, err := e.Client.DeleteIndex(req.Index).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &proto_config.DeleteIndexResponse{}, nil
}

// DeleteAlias from ES
func (e *elastic) DeleteAlias(ctx context.Context, req *proto_config.DeleteAliasRequest) (*proto_config.DeleteAliasResponse, error) {
	_, err := e.Client.Alias().Remove(req.Index, req.Alias).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &proto_config.DeleteAliasResponse{}, nil
}

// RenameAlias from ES
func (e *elastic) RenameAlias(ctx context.Context, req *proto_config.RenameAliasRequest) (*proto_config.RenameAliasResponse, error) {
	_, err := e.Client.Alias().Remove(req.Index, req.OldAlias).Do(ctx)
	if err != nil {
		return nil, err
	}

	_, err = e.AddAlias(ctx, &proto_config.AddAliasRequest{
		Index: req.Index,
		Alias: req.NewAlias,
	})
	if err != nil {
		return nil, err
	}

	return &proto_config.RenameAliasResponse{}, nil
}
