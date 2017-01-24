package tests

import (
	"bytes"
	"github.com/kazoup/platform/lib/globals"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var createindexdbtests = []struct {
	in  []byte
	out *http.Response
}{
	// Create same index twice
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_tests_index"
		}
	}`), &http.Response{StatusCode: 200}},
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.CreateIndex",
		"request": {
			"index": "db_config_srv_tests_index"
		}
	}`), &http.Response{StatusCode: 200}},
}

var addaliasdbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.AddAlias",
		"request": {
			"index": "db_config_srv_tests_index",
			"alias": "test_alias"
		}
	}`), &http.Response{StatusCode: 200}},
}

var renamealiasdbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.RenameAlias",
		"request": {
			"index": "db_config_srv_tests_index",
			"old_alias": "test_alias",
			"new_alias": "test_alias_renamed"
		}
	}`), &http.Response{StatusCode: 200}},
}

var deletealiasdbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteAlias",
		"request": {
			"index": "db_config_srv_tests_index",
			"alias": "test_alias_renamed"
		}
	}`), &http.Response{StatusCode: 200}},
}

var deleteindexdbtests = []struct {
	in  []byte
	out *http.Response
}{
	{[]byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.DeleteIndex",
		"request": {
			"index": "db_config_srv_tests_index"
		}
	}`), &http.Response{StatusCode: 200}},
}

func TestDBStatus(t *testing.T) {
	b := []byte(`{
		"service": "com.kazoup.srv.db",
		"method": "Config.Status",
		"request": {}
	}`)

	req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("Error create request %v", err)
	}

	req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
	req.Header.Add("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		t.Errorf("Error performing request with body: %s %v", string(b), err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != STATUS_OK {
		b, _ := ioutil.ReadAll(rsp.Body)
		t.Errorf("Expected %v with body %s, got %v", STATUS_OK, string(b), rsp.StatusCode, string(b))
	}
}

func TestDBCreateIndex(t *testing.T) {
	for _, v := range createindexdbtests {
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

func TestDBAddAlias(t *testing.T) {
	for _, v := range addaliasdbtests {
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

func TestDBRenameAlias(t *testing.T) {
	for _, v := range renamealiasdbtests {
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

func TestDBDeleteAlias(t *testing.T) {
	for _, v := range deletealiasdbtests {
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

func TestDBDeleteIndex(t *testing.T) {
	for _, v := range deleteindexdbtests {
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

/* THIS was used with local files,
func TestDBAggregate(t *testing.T) {

}*/
