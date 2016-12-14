package elastic

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/lib/globals"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	search_proto "github.com/kazoup/platform/search/srv/proto/search"
	lib "github.com/mattbaird/elastigo/lib"
	"log"
	"strconv"
)

func indexer(e *elastic) error {
	// Files
	go func() {
		for {
			select {
			case v := <-e.filesChannel:
				// File message can be notified, when a file is create, deleted or shared within kazoup
				if v.FileMessage.Notify {
					// We do not use bulk, as is just one element
					if _, err := e.conn.Index(v.FileMessage.Index, "file", v.FileMessage.Id, nil, v.FileMessage.Data); err != nil {
						log.Print("Bulk Indexer error %s", err)
					}

					n := &notification_proto.NotificationMessage{
						Method: globals.NOTIFY_REFRESH_SEARCH,
						UserId: v.FileMessage.UserId,
					}

					// Publish scan topic, crawlers should pick up message
					if err := v.Client.Publish(globals.NewSystemContext(), v.Client.NewPublication(globals.NotificationTopic, n)); err != nil {
						log.Print("Publishing (notify file) error %s", err)
					}
				} else {
					// Use bulk as we will index groups of documents
					if err := e.bulk.Index(v.FileMessage.Index, "file", v.FileMessage.Id, "", "", nil, v.FileMessage.Data); err != nil {
						log.Print("Bulk Indexer error %s", err)
					}
				}
			}

		}
	}()

	// Slack users
	go func() {
		for {
			select {
			case v := <-e.slackUsersChannel:
				if err := e.bulk.Index(v.Index, "user", v.Id, "", "", nil, v.Data); err != nil {
					log.Print("Bulk Indexer error %s", err)
				}
			}

		}
	}()

	// Slack channels
	go func() {
		for {
			select {
			case v := <-e.slackChannelsChannel:
				if err := e.bulk.Index(v.Index, "channel", v.Id, "", "", nil, v.Data); err != nil {
					log.Print("Bulk Indexer error %s", err)
				}
			}

		}
	}()

	return nil
}

func enricher(e *elastic) error {
	go func() {
		for {
			select {
			case v := <-e.crawlerFinished:
				log.Println(v)
			}
		}
	}()

	return nil
}

type JsonRemoveAliases struct {
	Actions []JsonAliasRemove `json:"actions"`
}

type JsonAliasRemove struct {
	Remove lib.JsonAlias `json:"remove"`
}

// The API allows you to remove an index alias through an API.
func (e *elastic) RemoveAlias(index string, alias string) (lib.BaseResponse, error) {
	var url string
	var retval lib.BaseResponse

	if len(index) > 0 {
		url = "/_aliases"
	} else {
		return retval, errors.New("alias required")
	}

	jsonAliases := JsonRemoveAliases{}
	jsonAliasRemove := JsonAliasRemove{}
	jsonAliasRemove.Remove.Alias = alias
	jsonAliasRemove.Remove.Index = index
	jsonAliases.Actions = append(jsonAliases.Actions, jsonAliasRemove)
	requestBody, err := json.Marshal(jsonAliases)

	if err != nil {
		return retval, err
	}

	body, err := e.conn.DoCommand("POST", url, nil, requestBody)
	if err != nil {
		return retval, err
	}

	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return retval, jsonErr
	}

	return retval, err
}

// TODO: use gabs (handle JSON in go)
// ElasticQuery to generate DSL query from params
type ElasticQuery struct {
	Index                string
	Id                   string
	UserId               string
	Term                 string
	From                 int64
	Size                 int64
	Category             string
	Url                  string
	Depth                int64
	Type                 string
	FileType             string
	LastSeen             int64
	NoKazoupFileOriginal bool
	Aggs                 []*search_proto.Aggregation
}

// Query generates a Elasticsearch DSL query
func (e *ElasticQuery) Query() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{`)
	buffer.WriteString(e.setSource())
	buffer.WriteString(e.filterFrom() + ",")
	buffer.WriteString(e.filterSize() + ",")
	buffer.WriteString(`"query": {"bool":{"must":[`)
	buffer.WriteString(e.queryTerm())
	buffer.WriteString(`], "filter":[`)
	buffer.WriteString(e.filterCategory() + ",")
	buffer.WriteString(e.filterDepth() + ",")
	buffer.WriteString(e.filterUrl() + ",")
	buffer.WriteString(e.filterUser() + ",")
	buffer.WriteString(e.filterLastSeen() + ",")
	buffer.WriteString(e.filterType())
	buffer.WriteString(`]}}, "sort":[`)
	buffer.WriteString(e.defaultSorting())
	buffer.WriteString(`]}`)

	return buffer.String(), nil
}

// Query generates a Elasticsearch DSL query
func (e *ElasticQuery) AggsQuery() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{"size":0,"query":{"filtered":{"filter":{"bool":{"must":[`)
	buffer.WriteString(e.filterType() + ",")
	buffer.WriteString(e.queryTerm() + ",")
	buffer.WriteString(e.filterCategory() + ",")
	buffer.WriteString(e.filterUrl())
	buffer.WriteString(`]}}}}, "aggs":{`)
	buffer.WriteString(e.aggs())
	buffer.WriteString(`}}`)

	return buffer.String(), nil
}

// Query generates a Elasticsearch DSL query
func (e *ElasticQuery) DeleteQuery() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{"query":`)
	buffer.WriteString(e.filterLastSeen())
	buffer.WriteString(`}`)

	return buffer.String(), nil
}

// QueryById generates a Elasticsearch DSL query for searching aliases by id
func (e *ElasticQuery) QueryById() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{"query":{"filtered":{"filter":{"bool":{"must":[{"term":{"id":"`)
	buffer.WriteString(e.Id + `"}}`)
	// Filter by user for files, not for users or channels (slack)
	// This is due to channels and users (slack) does not have to store the user they belong to
	if e.FileType == globals.FileType {
		buffer.WriteString(`,` + e.filterUser())
	}
	buffer.WriteString(`]}}}}}`)

	return buffer.String(), nil
}

func (e *ElasticQuery) defaultSorting() string {
	var buffer bytes.Buffer

	// Sorting for no datasource indexes
	if (e.From != 0 || e.Size != 0) && e.Index != globals.IndexDatasources && len(e.Term) == 0 {
		buffer.WriteString(`{"is_dir": "desc"},{"modified":"desc"},{"file_size": "desc"}`)
	}

	// Sorting for datasources
	if (e.From != 0 || e.Size != 0) && e.Index == globals.IndexDatasources {
		buffer.WriteString(`{"url": "desc"}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterLastSeen() string {
	var buffer bytes.Buffer

	if e.LastSeen > 0 {
		buffer.WriteString(`{"range":{"last_seen":{"lte":`)
		buffer.WriteString(strconv.Itoa(int(e.LastSeen)))
		buffer.WriteString(`}}}`)
	} else {
		buffer.WriteString(`{}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterType() string {
	var buffer bytes.Buffer

	switch e.FileType {
	case globals.FileTypeFile:
		buffer.WriteString(`{"term": {"is_dir": false}}`)
	case globals.FileTypeDirectory:
		buffer.WriteString(`{"term": {"is_dir": true}}`)
	default:
		buffer.WriteString(`{}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterDepth() string {
	var buffer bytes.Buffer

	if e.Depth <= 0 {
		buffer.WriteString(`{}`)
	} else {
		buffer.WriteString(`{"term": {"depth": "`)
		buffer.WriteString(strconv.FormatInt(e.Depth, 10))
		buffer.WriteString(`"}}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterUrl() string {
	var buffer bytes.Buffer

	if len(e.Url) <= 0 {
		buffer.WriteString(`{}`)
	} else {
		buffer.WriteString(`{"term": {"url": "`)
		buffer.WriteString(e.Url)
		buffer.WriteString(`"}}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterUser() string {
	var buffer bytes.Buffer

	// We filter by user_id when quering file or datasource ES documents.
	if len(e.UserId) > 0 && (e.Type == globals.FileType || e.Type == globals.TypeDatasource) {
		buffer.WriteString(`{"term": {"user_id": "`)
		buffer.WriteString(e.UserId)
		buffer.WriteString(`"}}`)
	} else {
		buffer.WriteString(`{}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) queryTerm() string {
	var buffer bytes.Buffer

	if len(e.Term) <= 0 {
		buffer.WriteString(`{}`)
	} else {
		buffer.WriteString(`{"match": {"name": "`)
		buffer.WriteString(e.Term)
		buffer.WriteString(`"}}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterCategory() string {
	var buffer bytes.Buffer

	if len(e.Category) <= 0 {
		buffer.WriteString(`{}`)
	} else {
		buffer.WriteString(`{"term": {"category": "`)
		buffer.WriteString(e.Category)
		buffer.WriteString(`"}}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterSize() string {
	var buffer bytes.Buffer

	buffer.WriteString(`"size": `)
	buffer.WriteString(strconv.FormatInt(e.Size, 10))

	return buffer.String()
}

func (e *ElasticQuery) filterFrom() string {
	var buffer bytes.Buffer

	buffer.WriteString(`"from": `)
	buffer.WriteString(strconv.FormatInt(e.From, 10))

	return buffer.String()
}

func (e *ElasticQuery) setSource() string {
	var buffer bytes.Buffer

	if e.NoKazoupFileOriginal {
		buffer.WriteString(`"_source": [
			"id",
			"user_id",
			"name",
			"url",
			"modified",
			"file_size",
			"is_dir",
			"category",
			"mime_type",
			"depth",
			"file_type",
			"last_seen",
			"datasource_id",
			"index"
		],`)
	}

	return buffer.String()
}

func (e *ElasticQuery) aggs() string {
	var buffer bytes.Buffer

	for k, v := range e.Aggs {
		buffer.WriteString(`"`)
		buffer.WriteString(strconv.Itoa(k))
		buffer.WriteString(`":{"`)
		buffer.WriteString(v.AggregationType)
		buffer.WriteString(`":{"field":"`)
		buffer.WriteString(v.Field)
		buffer.WriteString(`"}}`)

		if len(e.Aggs) > k+1 {
			buffer.WriteString(`,`)
		}
	}

	return buffer.String()
}
