package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"
)

var searcheable_set_data = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "search_index_test",
			"type": "file",
			"id": "test_search_id_1",
			"data": "{\"id\": \"test_search_id_1\", \"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"file_type\":\"type1\", \"access\":\"access1\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "search_index_test",
			"type": "file",
			"id": "test_search_id_2",
			"data": "{\"id\": \"test_search_id_2\", \"name\": \"orange\", \"category\": \"green\", \"last_seen\": 2, \"file_type\":\"type2\", \"access\":\"access1\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "search_index_test",
			"type": "file",
			"id": "test_search_id_3",
			"data": "{\"id\": \"test_search_id_3\", \"name\": \"lemon\", \"category\": \"green\", \"last_seen\": 3, \"file_type\":\"type3\", \"access\":\"access1\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "search_index_test",
			"type": "file",
			"id": "test_search_id_4",
			"data": "{\"id\": \"test_search_id_4\", \"name\": \"watermelon\", \"category\": \"black\", \"last_seen\": 4, \"file_type\":\"type1\", \"access\":\"access2\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "search_index_test",
			"type": "file",
			"id": "test_search_id_5",
			"data": "{\"id\": \"test_search_id_5\", \"name\": \"olive\", \"category\": \"black\", \"last_seen\": 5, \"file_type\":\"type2\", \"access\":\"access3\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var search_test = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "search_index_test",
			"from": 0,
			"size": 5
		}
	}`), &http.Response{StatusCode: 200}, time.Second},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "search_index_test",
			"from": 0,
			"size": 5,
			"term": "range"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "search_index_test",
			"from": 0,
			"size": 5,
			"category": "green"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "search_index_test",
			"from": 0,
			"size": 5,
			"file_type": "type1"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "search_index_test",
			"from": 0,
			"size": 5,
			"file_type": "access2"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

func TestSearchSearch(t *testing.T) {

	rangeTestTable(searcheable_set_data, t)

	rangeTestTableWithChecker(search_test, func(rsp *http.Response, t *testing.T) {
		type TestRsp struct {
			Info   string `json:"info"`
			Result string `json:"result"`
		}

		var tr TestRsp
		var tl map[string]int

		if err := json.NewDecoder(rsp.Body).Decode(&tr); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if err := json.Unmarshal([]byte(tr.Info), &tl); err != nil {
			t.Fatalf("Error unmarshalling response: %v", err)
		}

		log.Println(tl["total"])

	}, t)
}
