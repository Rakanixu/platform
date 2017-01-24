package tests

import (
	"bytes"
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var createdbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 1}"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Create",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"id": "test_id_2",
			"data": "{\"name\": \"grass\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"string\",\"bool\": true,\"int\": 2}"
		}
	}`), &http.Response{StatusCode: 200}},
}

var readdbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Read",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"id": "test_id"
		}
	}`), &http.Response{StatusCode: 200}},
}

var updatedbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Update",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"id": "test_id_1",
			"data": "{\"name\": \"tree\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"updated string\",\"bool\": false,\"int\": 0}"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Update",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"id": "test_id_2",
			"data": "{\"name\": \"grass\", \"category\": \"green\", \"last_seen\": 1, \"string\": \"updated string\",\"bool\": false,\"int\": 0}"
		}
	}`), &http.Response{StatusCode: 200}},
}

var searchdbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Search",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"user_id": "` + USER_ID + `",
			"term": "tree"
		}
	}`), &http.Response{StatusCode: 200}},
}

/*var searchbyiddbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.SearchById",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"id": "test_id",
			"user_id": "` + USER_ID + `",
			"name": "tree"
		}
	}`), &http.Response{StatusCode: 200}},
}*/

var deletedbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.Delete",
		"request": {
			"index": "db_srv_tests_index",
			"type": "test_document",
			"id": "test_id_1"
		}
	}`), &http.Response{StatusCode: 200}},
}

var deletebyquerydbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "DB.DeleteByQuery",
		"request": {
			"indexes": ["db_srv_tests_index"],
			"types": ["test_document"],
			"last_seen": 1
		}
	}`), &http.Response{StatusCode: 200}},
}

func TestDBCreate(t *testing.T) {
	for _, v := range createdbtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != v.out.StatusCode {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}
	}

	time.Sleep(time.Second)
}

func TestDBRead(t *testing.T) {
	for _, v := range readdbtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != v.out.StatusCode {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}
	}
}

func TestDBUpdate(t *testing.T) {
	for _, v := range updatedbtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != v.out.StatusCode {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}
	}
	time.Sleep(time.Second)
}

func TestDBSearch(t *testing.T) {
	for _, v := range searchdbtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != v.out.StatusCode {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Fatalf("Expected %v with body %s, got %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}

		type TestRsp struct {
			Result string `json:"result"`
			Info   string `json:"info"`
		}

		var tr TestRsp
		var tl map[string]int

		if err := json.NewDecoder(rsp.Body).Decode(&tr); err != nil {
			t.Errorf("Error decoding response: %v", err)
		}

		if err := json.Unmarshal([]byte(tr.Info), &tl); err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
		}

		if tl["total"] != 1 {
			t.Errorf("Expected 1 result, got ", tl["total"])
		}
	}
}

// Requires to add user_id to prototype
/*
func TestDBSearchById(t *testing.T) {
	for _, v := range searchbyiddbtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != v.out.StatusCode {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}
	}
}
*/

func TestDBDelete(t *testing.T) {
	for _, v := range deletedbtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != v.out.StatusCode {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}
	}
}

func TestDBDeleteByQuery(t *testing.T) {
	for _, v := range deletebyquerydbtests {
		req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(v.in))
		if err != nil {
			t.Errorf("Error create request %v", err)
		}

		req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
		req.Header.Add("Content-Type", "application/json")
		rsp, err := c.Do(req)
		if err != nil {
			t.Errorf("Error performing request with body: %s %v", string(v.in), err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != v.out.StatusCode {
			b, _ := ioutil.ReadAll(rsp.Body)
			t.Errorf("Expected %v with body %s, got %v", v.out.StatusCode, string(v.in), rsp.StatusCode, string(b))
		}
	}
}
