package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

var db_create = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_create_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_create_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_read = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_read_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Read",
		"request": {
			"index": "db_srv_read_test",
			"type": "test_document",
			"id": "test_id"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_read_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_update = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_update_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Update",
		"request": {
			"index": "db_srv_update_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"updated string\",\"bool\": false,\"int\": 0}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_update_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_delete = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_delete_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Delete",
		"request": {
			"index": "db_srv_delete_test",
			"type": "test_document",
			"id": "test_id_1"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_delete_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_deletebyquery = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_deletebyquery_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.DeleteByQuery",
		"request": {
			"indexes": ["db_srv_deletebyquery_test"],
			"types": ["test_document"],
			"last_seen": 2
		}
	}`), &http.Response{StatusCode: 200}, time.Second},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_deletebyquery_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_create_for_search = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_search_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_search_test",
			"type": "test_document",
			"id": "test_id_2",
			"data": "{\"name\": \"orange\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_search = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Search",
		"request": {
			"index": "db_srv_search_test",
			"type": "test_document",
			"user_id": "` + USER_ID + `",
			"term": "tree"
		}
	}`), &http.Response{StatusCode: 200}, time.Second},
}

var db_deleteindex_for_search = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_search_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_create_for_searchbyid = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_searchbyid_test",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"id\": \"test_id_1\", \"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_searchbyid_test",
			"type": "test_document",
			"id": "test_id_2",
			"data": "{\"id\": \"test_id_2\", \"name\": \"orange\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

var db_searchbyid = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.SearchById",
		"request": {
			"index": "db_srv_searchbyid_test",
			"type": "test_document",
			"id": "test_id_1",
			"user_id": "` + USER_ID + `",
			"name": "tree"
		}
	}`), &http.Response{StatusCode: 200}, time.Second},
}
var db_deleteindex_for_searchbyid = testTable{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_srv_searchbyid_test"
		}
	}`), &http.Response{StatusCode: 200}, noDuration},
}

func TestDBCreate(t *testing.T) {
	// Create document, delete index
	rangeTestTable(db_create, JWT_TOKEN_USER_1, t)
}

func TestDBRead(t *testing.T) {
	// Create document, read document, delete index
	rangeTestTable(db_read, JWT_TOKEN_USER_1, t)
}

func TestDBUpdate(t *testing.T) {
	// Create document, update document, delete index
	rangeTestTable(db_update, JWT_TOKEN_USER_1, t)
}

func TestDBDelete(t *testing.T) {
	// Create document, delete document, delete index
	rangeTestTable(db_delete, JWT_TOKEN_USER_1, t)
}

func TestDBDeleteByQuery(t *testing.T) {
	// Create document, delete document, delete index
	rangeTestTable(db_deletebyquery, JWT_TOKEN_USER_1, t)
}

func TestDBSearch(t *testing.T) {
	// Create 2 document
	rangeTestTable(db_create_for_search, JWT_TOKEN_USER_1, t)

	// Search by term for 1 document
	rangeTestTableWithChecker(db_search, JWT_TOKEN_USER_1, func(rsp *http.Response, t *testing.T) {
		type TestRsp struct {
			Result string `json:"result"`
			Info   string `json:"info"`
		}

		var tr TestRsp
		var tl map[string]int

		if err := json.NewDecoder(rsp.Body).Decode(&tr); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if err := json.Unmarshal([]byte(tr.Info), &tl); err != nil {
			t.Fatalf("Error unmarshalling response: %v", err)
		}

		if tl["total"] != 1 {
			t.Errorf("Expected 1 result, got %v", tl["total"])
		}
	}, t)

	// Delete index used for test
	rangeTestTable(db_deleteindex_for_search, JWT_TOKEN_USER_1, t)
}

// Requires to add user_id to prototype
func TestDBSearchById(t *testing.T) {
	// Create 2 document
	rangeTestTable(db_create_for_searchbyid, JWT_TOKEN_USER_1, t)

	// Search by id for 1 document
	rangeTestTableWithChecker(db_searchbyid, JWT_TOKEN_USER_1, func(rsp *http.Response, t *testing.T) {
		type TestRsp struct {
			Result string `json:"result"`
		}

		var tr TestRsp

		if err := json.NewDecoder(rsp.Body).Decode(&tr); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		// Search by id returns an object stringify, so empty result will be "{}"
		if len(tr.Result) <= 2 {
			t.Errorf("Expected string with result, got: %v", tr.Result)
		}
	}, t)

	// Delete index used for test
	rangeTestTable(db_deleteindex_for_searchbyid, JWT_TOKEN_USER_1, t)
}
