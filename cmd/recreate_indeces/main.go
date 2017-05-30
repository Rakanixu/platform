package main

/*

Usage: Apply new settings and mapping, afterwars run this tool.
It will remove all indeces that follows name pattern "index*"
New indeces will be created, and its old aliases will be applied
*/

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	url := flag.String("url", "http://localhost:9200", "Elastic endpoint")
	//url := flag.String("url", "https://admin:9svne8f655h@8d22518314b8a6bab84906817730e7f4.eu-west-1.aws.found.io:9243", "Elastic endpoint")

	flag.Parse()

	// Retrieve all indeces with its aliases
	rc, err := doRequest(http.MethodGet, fmt.Sprintf("%s/_aliases", *url), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	var srsp map[string]map[string]interface{}
	b, err := ioutil.ReadAll(rc)
	if err := json.Unmarshal(b, &srsp); err != nil {
		log.Fatal(err)
	}

	// We match our indeces that starts by index
	re := regexp.MustCompile("index")
	for k, v := range srsp {
		if re.MatchString(k) {
			// Delete the index
			err := logRequest(http.MethodDelete, fmt.Sprintf("%s/%s", *url, k), nil)
			if err != nil {
				log.Fatal(err)
			}

			var i int32
			var bf bytes.Buffer
			bf.WriteString(`{"aliases" : {`)

			for alias, _ := range v["aliases"].(map[string]interface{}) {
				if i != 0 {
					bf.WriteString(`,`)
				}
				bf.WriteString(`"` + alias + `" : {}`)
				i++
			}
			bf.WriteString(`}}`)

			// Recreate index with all its alias
			err = logRequest(http.MethodPut, fmt.Sprintf("%s/%s", *url, k), bytes.NewReader(bf.Bytes()))
			if err != nil {
				log.Fatal(err)
			}
		}
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
