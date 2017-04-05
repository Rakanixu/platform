package main

// ElasticSearch v5.X
// This tool finds all indeces in an elasticsearch cluster to update its settings and mapping
// Check update mapping constrains:
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-mapping.html#updating-field-mappings
import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	INDECES_URL = "/*/_stats/store"
)

func main() {
	fa := flag.String("a", "analysis.json", "Path to ES new analysis json file")
	fm := flag.String("f", "map.json", "Path to ES new mapping json file")
	url := flag.String("url", "http://elasticsearch:9200", "Path to ES map file")

	flag.Parse()
	if *fa == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *fm == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rc, err := doRequest(http.MethodGet, fmt.Sprintf("%s%s", *url, INDECES_URL), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	var srsp map[string]map[string]interface{}
	b, err := ioutil.ReadAll(rc)
	if err := json.Unmarshal(b, &srsp); err != nil {
		log.Fatal(err)
	}

	var bf bytes.Buffer
	bf.WriteString("/")

	// We match our indeces that starts by index
	re := regexp.MustCompile("index")
	for k, _ := range srsp["indices"] {
		if re.MatchString(k) {
			// Concatenate all indeces
			bf.WriteString(k)
			bf.WriteString(",")

			// Close index
			err := logRequest(http.MethodPost, fmt.Sprintf("%s/%s/_close", *url, k), nil)
			if err != nil {
				log.Fatal(err)
			}

			// File will be open on every iteration
			ea, err := os.Open(*fa)
			if err != nil {
				log.Fatal(err)
			}

			// Update analysis
			err = logRequest(http.MethodPut, fmt.Sprintf("%s/%s/_settings", *url, k), ea)
			if err != nil {
				log.Fatal(err)
			}

			// Open index
			err = logRequest(http.MethodPost, fmt.Sprintf("%s/%s/_open", *url, k), nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	em, err := os.Open(*fm)
	if err != nil {
		log.Fatal(err)
	}

	err = logRequest(http.MethodPut, fmt.Sprintf("%s%s/_mapping/%s", *url, strings.TrimSuffix(bf.String(), ","), globals.FileType), em)
	if err != nil {
		log.Fatal(err)
	}
}

func logRequest(method, url string, body io.Reader) error {
	rc, err := doRequest(method, url, body)
	if err != nil {
		return err
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	log.Println(url, string(b))

	return nil
}

func doRequest(method, url string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}
