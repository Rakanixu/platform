package elastic

import (
	"bytes"
	"github.com/kazoup/platform/lib/globals"
	"strconv"
)

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
	Access               string
	NoKazoupFileOriginal bool
}

// Query generates a Elasticsearch DSL query
func (e *ElasticQuery) Query() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{`)
	buffer.WriteString(e.setSource())
	buffer.WriteString(e.filterFrom() + ",")
	buffer.WriteString(e.filterSize() + ",")
	buffer.WriteString(`"query": {"bool":{"must":{`)
	buffer.WriteString(`"bool":{"should":[`)
	buffer.WriteString(e.queryTerm() + `,`)
	buffer.WriteString(e.queryContent() + `,`)
	buffer.WriteString(e.queryTags())
	buffer.WriteString(`]}`)
	buffer.WriteString(`}, "filter":[`)
	buffer.WriteString(e.filterCategory() + ",")
	buffer.WriteString(e.filterDepth() + ",")
	buffer.WriteString(e.filterUrl() + ",")
	buffer.WriteString(e.filterUser() + ",")
	buffer.WriteString(e.filterLastSeen() + ",")
	buffer.WriteString(e.filterType() + ",")
	buffer.WriteString(e.filterAccess())
	buffer.WriteString(`]}}, "sort":[`)
	buffer.WriteString(e.defaultSorting())
	buffer.WriteString(`]`)
	buffer.WriteString(e.contentHighlight())
	buffer.WriteString(`}`)

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

	buffer.WriteString(`{"query":{"bool":{"must":[{"term":{"id":"`)
	buffer.WriteString(e.Id + `"}}`)
	// Filter by user for files, not for users or channels (slack)
	// This is due to channels and users (slack) does not have to store the user they belong to
	if e.FileType == globals.FileType {
		buffer.WriteString(`,` + e.filterUser())
	}
	buffer.WriteString(`]}}}`)

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

func (e *ElasticQuery) filterAccess() string {
	var buffer bytes.Buffer

	if len(e.Access) > 0 {
		buffer.WriteString(`{"term": {"access": "` + e.Access + `"}}`)
	} else {
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

	buffer.WriteString(`{"term": {"user_id": "`)
	buffer.WriteString(e.UserId)
	buffer.WriteString(`"}}`)

	return buffer.String()
}

func (e *ElasticQuery) queryTerm() string {
	var buffer bytes.Buffer

	if len(e.Term) <= 0 {
		buffer.WriteString(`{}`)
	} else {
		buffer.WriteString(`{"match": {"name.raw":{"boost":10,"query": "`)
		buffer.WriteString(e.Term)
		buffer.WriteString(`"}}},`)
		buffer.WriteString(`{"match": {"name":{"query": "`)
		buffer.WriteString(e.Term)
		buffer.WriteString(`"}}}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) queryContent() string {
	var buffer bytes.Buffer

	if len(e.Term) > 0 && e.Type == globals.FileType {
		buffer.WriteString(`{"match_phrase": {"content":{"boost":6,"query":"`)
		buffer.WriteString(e.Term)
		buffer.WriteString(`"}}}`)
	} else {
		buffer.WriteString(`{}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) queryTags() string {
	var buffer bytes.Buffer

	if len(e.Term) > 0 && e.Type == globals.FileType {
		buffer.WriteString(`{"match": {"tags":{"boost":6,"query":"`)
		buffer.WriteString(e.Term)
		buffer.WriteString(`"}}}`)
	} else {
		buffer.WriteString(`{}`)
	}

	return buffer.String()
}

func (e *ElasticQuery) contentHighlight() string {
	var buffer bytes.Buffer

	if len(e.Term) > 0 && e.Type == globals.FileType {
		buffer.WriteString(`,"highlight":{"fields":{"content":{"number_of_fragments": 1,"fragment_size":150}}}`)
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
			"access",
			"content_category",
			"datasource_id",
			"index"
		],`)
	}

	return buffer.String()
}
