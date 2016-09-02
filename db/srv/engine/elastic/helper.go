package elastic

import (
	"bytes"
	"encoding/json"
	"errors"
	lib "github.com/mattbaird/elastigo/lib"
	"log"
	"strconv"
)

func indexer(e *elastic) error {
	go func() {
		for {
			select {
			case v := <-e.filesChannel:
				if err := e.bulk.Index(v.Index, "file", v.Id, "", "", nil, v.Data); err != nil {
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

// TODO: use gabs (handle JSON in go)
// ElasticQuery to generate DSL query from params
type ElasticQuery struct {
	Term     string
	From     int64
	Size     int64
	Category string
	Url      string
	Depth    int64
	Type     string
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

// Query generates a Elasticsearch DSL query
func (e *ElasticQuery) Query() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{`)
	buffer.WriteString(e.filterFrom() + ",")
	buffer.WriteString(e.filterSize() + ",")
	buffer.WriteString(`"query": {"bool":{"must":[`)
	buffer.WriteString(e.queryTerm())
	buffer.WriteString(`], "filter":[`)
	buffer.WriteString(e.filterCategory() + ",")
	buffer.WriteString(e.filterDepth() + ",")
	buffer.WriteString(e.filterUrl() + ",")
	buffer.WriteString(e.filterType())
	buffer.WriteString(`]}}, "sort":[`)
	buffer.WriteString(e.defaultSorting())
	buffer.WriteString(`]}`)

	return buffer.String(), nil
}

func (e *ElasticQuery) defaultSorting() string {
	var buffer bytes.Buffer

	if (e.From != 0 || e.Size != 0) && e.Type == "file" {
		buffer.WriteString(`{"is_dir": "desc"},{"size": "desc"}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) filterType() string {
	var buffer bytes.Buffer

	switch e.Type {
	case "files":
		buffer.WriteString(`{"term": {"is_dir": false}}`)
	case "directories":
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
