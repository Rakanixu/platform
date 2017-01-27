package tests

import (
	"bytes"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	RPC_ENPOINT = "https://web.kazoup.io:8082/rpc"
	USER_ID     = "test@kazoup.com"
	USER_PWD    = "ksu4awemtest"
	STATUS_OK   = 200
)

const noDuration time.Duration = 0

const (
	BOX_URL          = "box://" + USER_ID
	DROPBOX_URL      = "dropbox://" + USER_ID
	GOOGLE_DRIVE_URL = "googledrive://" + USER_ID
	ONE_DRIVE_URL    = "onedrive://" + USER_ID
	GMAIL_URL        = "gmail://" + USER_ID
	SLACK_URL        = "slack://" + USER_ID
)

var datasources []proto.Endpoint

type testTable []struct {
	in    []byte
	out   *http.Response
	delay time.Duration // Timeout before request is done
}

type Checker func(*http.Response, *testing.T)

var c *http.Client

func init() {
	c = &http.Client{}
}

func makeRequest(body []byte, result *http.Response, t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error create request %v", err)
	}

	req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
	req.Header.Add("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Error performing request with body: %s %v", string(body), err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != result.StatusCode {
		b, _ := ioutil.ReadAll(rsp.Body)
		t.Errorf("Expected %v with body %s, got %v", result.StatusCode, string(body), rsp.StatusCode, string(b))
	}
}

func rangeTestTable(tt testTable, t *testing.T) {
	for _, v := range tt {
		time.Sleep(v.delay)
		makeRequest(v.in, v.out, t)
	}
}

func makeRequestWithChecker(body []byte, result *http.Response, ch Checker, t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error create request %v", err)
	}

	req.Header.Add("Authorization", globals.SYSTEM_TOKEN)
	req.Header.Add("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Error performing request with body: %s %v", string(body), err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != result.StatusCode {
		b, _ := ioutil.ReadAll(rsp.Body)
		t.Errorf("Expected %v with body %s, got %v", result.StatusCode, string(body), rsp.StatusCode, string(b))
	}

	ch(rsp, t)
}

func rangeTestTableWithChecker(tt testTable, ch Checker, t *testing.T) {
	for _, v := range tt {
		time.Sleep(v.delay)
		makeRequestWithChecker(v.in, v.out, ch, t)
	}
}
