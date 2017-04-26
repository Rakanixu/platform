package elastic

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/custom/proto/custom"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
	"io"
	"os"
)

type elastic struct {
	Client *elib.Client
}

func init() {
	custom.Register(new(elastic))
}

// Init config db
func (e *elastic) Init() error {
	var err error

	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
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

// ScrollUnprocessedAudio retrieves audio files not processed yet
func (e *elastic) ScrollUnprocessedFiles(ctx context.Context, req *proto_custom.ScrollUnprocessedFilesRequest) (*proto_custom.ScrollUnprocessedFilesResponse, error) {
	var results []interface{}
	var err error
	var rstr string
	var uID string
	done := false

	// Get user id from context
	uID, err = globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	eQuery := ElasticQuery{
		Index:    req.Index,
		Field:    req.Field,
		UserId:   uID,
		Category: req.Category,
	}
	query, err := eQuery.ScrollUnprocessedFile()
	if err != nil {
		return nil, err
	}

	s := e.Client.Scroll(req.Index).Type(globals.FileType).Body(query)
	out, err := s.Do(ctx)
	if err == io.EOF {
		done = true
	}
	if err != io.EOF && err != nil {
		return nil, err
	}

	results, err = attachResults(results, out.Hits)
	if err != nil {
		return nil, err
	}

	if !done {
		results, err = scroll(results, s, out.ScrollId)
		if err != nil {
			return nil, err
		}
	}

	if len(results) == 0 {
		rstr = `[]`
	} else {
		b, err := json.Marshal(results)
		if err != nil {
			return nil, err
		}

		rstr = string(b)
	}

	return &proto_custom.ScrollUnprocessedFilesResponse{
		Result: rstr,
	}, nil
}
