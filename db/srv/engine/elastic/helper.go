package elastic

import (
	"bytes"
	"log"
	"strconv"
)

func indexer(e *elastic) error {
	go func() {
		for {
			select {
			case v := <-e.filesChannel:
				if err := e.bulk.Index("files", "file", v.Id, "", "", nil, v.Data); err != nil {
					log.Print("Bulk Indexer error %s", err)
				}
			}

		}
	}()

	return nil
}

type ElasticQuery struct {
	Term     string
	From     int64
	Size     int64
	Category string
	Url      string
	Depth    int64
	Type     string
}

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

	log.Println(buffer.String())

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
