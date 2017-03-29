package main

// ElasticSearch v5.X
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	INDECES_URL = "/*/_stats/store"
)

func main() {
	f := flag.String("f", "map.json", "Path to ES map file")
	url := flag.String("url", "http://elasticsearch:9200", "Path to ES map file")
	flag.Parse()
	if *f == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	_, err := os.Open(*f)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", *url, INDECES_URL), nil)
	if err != nil {
		log.Fatal(err)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	var srsp map[string]map[string]interface{}

	b, err := ioutil.ReadAll(rsp.Body)

	if err := json.Unmarshal(b, &srsp); err != nil {
		log.Fatal(err)
	}

	// We are lloking for the keys, that are the index name in ES
	for k, _ := range srsp["indices"] {
		log.Println(k)
	}

}
