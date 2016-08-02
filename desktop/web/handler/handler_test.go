package handler_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kazoup/platform/desktop/web/handler"
)

var Tests = []struct {
	path   string
	status int
}{
	{"", http.StatusBadRequest},
	{"dhsdhjsdh", http.StatusBadRequest},
	{"?path=/home/file/path", http.StatusOK},
}

func TestResponseCodes(t *testing.T) {
	ts := httptest.NewServer(handler.NewStreamHandler("."))
	defer ts.Close()

	for _, test := range Tests {

		res, err := http.Get(ts.URL + test.path)
		if err != nil {
			log.Fatal(err)
		}
		_, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		//Check if we getting http.StatusBadRequest
		if res.StatusCode != test.status {
			t.Errorf("Expected status code %s got %s with path %s", test.status, res.StatusCode, test.path)
		}
	}
}
