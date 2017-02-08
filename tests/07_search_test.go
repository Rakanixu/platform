package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

var searcheable_set_data = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "index_search_index_test",
			"type": "file",
			"id": "test_search_id_1",
			"data": "{\"id\": \"test_search_id_1\", \"user_id\": \"` + USER_ID + `\", \"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"is_dir\":false, \"file_size\":1000, \"modified\":\"016-11-04T13:21:24Z\", \"file_type\":\"files\", \"access\":\"access1\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "index_search_index_test",
			"type": "file",
			"id": "test_search_id_2",
			"data": "{\"id\": \"test_search_id_2\", \"user_id\": \"` + USER_ID + `\", \"name\": \"orange\", \"category\": \"green\", \"last_seen\": 2, \"is_dir\":false, \"file_size\":1000, \"modified\":\"016-11-04T13:21:24Z\", \"file_type\":\"files\", \"access\":\"access1\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "index_search_index_test",
			"type": "file",
			"id": "test_search_id_3",
			"data": "{\"id\": \"test_search_id_3\", \"user_id\": \"` + USER_ID + `\", \"name\": \"lemon\", \"category\": \"green\", \"last_seen\": 3, \"is_dir\":false, \"file_size\":1000, \"modified\":\"016-11-04T13:21:24Z\", \"file_type\":\"files\", \"access\":\"access1\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "index_search_index_test",
			"type": "file",
			"id": "test_search_id_4",
			"data": "{\"id\": \"test_search_id_4\", \"user_id\": \"` + USER_ID + `\", \"name\": \"watermelon\", \"category\": \"black\", \"last_seen\": 4, \"is_dir\":false, \"file_size\":1000, \"modified\":\"016-11-04T13:21:24Z\", \"file_type\":\"files\", \"access\":\"access2\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "index_search_index_test",
			"type": "file",
			"id": "test_search_id_5",
			"data": "{\"id\": \"test_search_id_5\", \"user_id\": \"` + USER_ID + `\", \"name\": \"olive\", \"category\": \"black\", \"last_seen\": 5, \"is_dir\":false, \"file_size\":1000, \"modified\":\"016-11-04T13:21:24Z\", \"file_type\":\"files\", \"access\":\"access3\"}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var search_test_count_1 = []int{5, 1, 0, 3, 1}
var search_test_count_2 = []int{0, 0, 0, 0, 0}
var search_test = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "index_search_index_test",
			"from": 0,
			"size": 5
		}
	}`), &http.Response{StatusCode: 200}, time.Second},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "index_search_index_test",
			"from": 0,
			"size": 5,
			"term": "range"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "index_search_index_test",
			"from": 0,
			"size": 5,
			"term": "www"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "index_search_index_test",
			"from": 0,
			"size": 5,
			"category": "green"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.search",
		"method": "Search.Search",
		"request": {
			"index": "index_search_index_test",
			"from": 0,
			"size": 5,
			"access": "access2"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var search_deleteindex = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "index_search_index_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

func TestSearchSearch(t *testing.T) {
	i := 0

	rangeTestTable(searcheable_set_data, JWT_TOKEN_USER_1, dbAccessHeader, t)

	rangeTestTableWithChecker(search_test, JWT_TOKEN_USER_1, emptyHeader, func(rsp *http.Response, t *testing.T) {
		defer rsp.Body.Close()
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

		if tl["total"] != search_test_count_1[i] {
			t.Errorf("Expecting %v results, got %v", search_test_count_1[i], tl["total"])
		}

		i++
	}, t)

	rangeTestTable(search_deleteindex, JWT_TOKEN_USER_1, emptyHeader, t)
}

func TestSearchSearchScopedData(t *testing.T) {
	i := 0

	rangeTestTable(searcheable_set_data, JWT_TOKEN_USER_2, dbAccessHeader, t)

	rangeTestTableWithChecker(search_test, JWT_TOKEN_USER_2, emptyHeader, func(rsp *http.Response, t *testing.T) {
		defer rsp.Body.Close()
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

		if tl["total"] != search_test_count_2[i] {
			t.Errorf("Expecting %v results, got %v", search_test_count_2[i], tl["total"])
		}

		i++
	}, t)

	rangeTestTable(search_deleteindex, JWT_TOKEN_USER_2, emptyHeader, t)
}
