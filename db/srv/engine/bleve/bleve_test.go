package bleve

import (
	lib "github.com/blevesearch/bleve"
	"os"
	"testing"
)

var queries = []struct {
	in   string
	hits uint64
}{
	{"hello", 5},
	{"bye", 1},
	{"+category:bye", 1},
}

var data = []struct {
	Id       string
	Category string `json:"category"`
}{
	{"1", "hello"},
	{"2", "hello"},
	{"3", "hello"},
	{"4", "hello"},
	{"5", "hello"},
	{"6", "bye"},
	{"7", "byee"},
}

func TestBleve_Search(t *testing.T) {
	files, err := lib.Open("/tmp/test/idx")
	if err != nil {
		mapping := lib.NewIndexMapping()
		//mapping.DefaultAnalyzer = keyword_analyzer.Name

		files, err = lib.New("/tmp/test/idx", mapping)
		if err != nil {
			t.Errorf("%v", err)
		}
	}

	for _, v := range data {
		files.Index(v.Id, v)
	}

	for _, v := range queries {
		qString := v.in
		q := lib.NewQueryStringQuery(qString)
		sr := lib.NewSearchRequestOptions(q, 6, 0, false)
		sr.Fields = []string{"*"} // Retrieve all fields

		results, err := files.Search(sr)
		if err != nil {
			t.Errorf("%v", err)
		}

		if v.hits != results.Total {
			t.Errorf("Expecting %d but got %d", v.hits, results.Total)
		}

	}

	if err := os.RemoveAll("/tmp/test"); err != nil {
		t.Fatalf("%v", err)
	}

}
