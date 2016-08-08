package elasticquery

import (
	"bytes"
	"github.com/kazoup/platform/search/srv/query"
	"log"
	"strconv"
)

type ElasticQuery struct {
	Term     string
	From     int64
	Size     int64
	Category string
	Querier  query.Querier
}

func (e *ElasticQuery) Query() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{`)
	buffer.WriteString(e.filterFrom() + ",")
	buffer.WriteString(e.filterSize() + ",")
	buffer.WriteString(`"query": {"bool":{"must":[`)
	buffer.WriteString(e.filterTerm())
	buffer.WriteString(`], "filter":[`)
	buffer.WriteString(e.filterCategory())
	buffer.WriteString(`]}}}`)

	log.Println(buffer.String())

	return buffer.String(), nil
}

func (e *ElasticQuery) filterTerm() string {
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

	buffer.WriteString(`"size": "`)
	buffer.WriteString(strconv.FormatInt(e.Size, 10))
	buffer.WriteString(`"`)

	return buffer.String()
}

func (e *ElasticQuery) filterFrom() string {
	var buffer bytes.Buffer

	buffer.WriteString(`"from": "`)
	buffer.WriteString(strconv.FormatInt(e.From, 10))
	buffer.WriteString(`"`)

	return buffer.String()
}
