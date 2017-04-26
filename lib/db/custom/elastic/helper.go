package elastic

import (
	"bytes"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/normalization/text"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
	"io"
	"strings"
)

// ElasticQuery to generate DSL query from params
type ElasticQuery struct {
	Index    string
	UserId   string
	Category string
	Field    string
}

// ScrollUnprocessedFile generates a query to retrieve files not processed before they were modified
func (e *ElasticQuery) ScrollUnprocessedFile() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(`{"_source": ["id", "file_type", "index"],`)
	buffer.WriteString(`"size": 1000,`)
	buffer.WriteString(`"query": {"bool":{"must":{`)
	buffer.WriteString(`"bool":{"should":[`)
	buffer.WriteString(`{
	     "bool":{
		"must_not":{
		   "exists":{
		      "field":"` + e.Field + `"
		   }
		}
	     }
	  },
	  {
	     "bool":{
		"must":[
		   {
		      "exists":{
			 "field":"` + e.Field + `"
		      }
		   },
		   {
		      "script":{
			 "script":{
			    "inline":"doc['` + e.Field + `'].value < doc['modified'].value",
			    "lang":"painless"
			 }
		      }
		   }
		]
	     }
	  }`)
	buffer.WriteString(`]}`)
	buffer.WriteString(`}, "filter":[`)
	buffer.WriteString(e.filterCategory() + ",")
	buffer.WriteString(e.filterUser())
	buffer.WriteString(`]}}}`)

	q, err := text.ReplaceTabs(buffer.String())
	if err != nil {
		return "", err
	}

	return text.ReplaceNewLines(strings.Replace(q, " ", "", -1))
}

func (e *ElasticQuery) filterUser() string {
	var buffer bytes.Buffer

	buffer.WriteString(`{"term": {"user_id": "`)
	buffer.WriteString(e.UserId)
	buffer.WriteString(`"}}`)

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

// scroll recursively paginate over scroll service until retrieves all documents
func scroll(results []interface{}, scrollSrv *elib.ScrollService, scrollId string) ([]interface{}, error) {
	done := false
	out, err := scrollSrv.ScrollId(scrollId).Do(context.Background())
	if err == io.EOF {
		done = true

		return results, nil
	}
	if err != io.EOF && err != nil {
		return nil, err
	}

	results, err = attachResults(results, out.Hits)
	if err != nil {
		return nil, err
	}

	if !done {
		return scroll(results, scrollSrv, out.ScrollId)
	}

	return results, nil
}

// attachResults helps to appends matched documents
func attachResults(results []interface{}, hits *elib.SearchHits) ([]interface{}, error) {
	for _, v := range hits.Hits {
		data, err := v.Source.MarshalJSON()
		if err != nil {
			return nil, err
		}

		f, err := file.NewFileFromString(string(data))

		if err != nil {
			return nil, err
		}

		results = append(results, f)
	}

	return results, nil
}
