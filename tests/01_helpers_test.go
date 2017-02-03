package tests

import (
	"bytes"
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

const (
	RPC_ENPOINT = "https://web.kazoup.io:8082/rpc"
	USER_NAME   = "test@kazoup.io"
	USER_NAME_2 = "test2@kazoup.io"
	USER_ID     = "auth0|58944e290861f87f6329526d"
	USER_PWD    = "ksu4awemtest"
	STATUS_OK   = 200
)

const noDuration time.Duration = 0

const (
	BOX_URL          = "box://" + USER_NAME
	DROPBOX_URL      = "dropbox://" + USER_NAME
	GOOGLE_DRIVE_URL = "googledrive://" + USER_NAME
	ONE_DRIVE_URL    = "onedrive://" + USER_NAME
	GMAIL_URL        = "gmail://" + USER_NAME
	SLACK_URL        = "slack://" + USER_NAME
)

var datasources []proto.Endpoint
var (
	JWT_TOKEN_USER_1 string
	JWT_TOKEN_USER_2 string
	JWT_INVALID      string
)

type testTable []struct {
	in    []byte
	out   *http.Response
	delay time.Duration // Timeout before request is done
}

type Checker func(*http.Response, *testing.T)

type authRsp struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

var c *http.Client

func init() {
	var err error

	c = &http.Client{}

	JWT_TOKEN_USER_1, err = authenticateTestUser(USER_NAME, USER_PWD)
	if err != nil {
		panic(err)
	}

	JWT_TOKEN_USER_2, err = authenticateTestUser(USER_NAME_2, USER_PWD)
	if err != nil {
		panic(err)
	}

	JWT_INVALID = "randomstringwithwithbitsofvalidonef0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2them91cC5ldS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NTg5NDRlMjkwODYxZjg3ZjYzMjk1MjZkIiwiYXVkIjoiNnpJRG04SW5oYlRScDFiTDJDNG0xVEs0TGxyNGFyVHkiLCJleHAiOjE0ODYxNTQ4OTQsImlhdCI6MTQ4NjExODg5NCwiYXpwIjoiNU9DSll1VHE1RG9nOTYwYzNsZlZFc0JscXVEWDlLYTIifQ.M9PvX8kErBeC2In5JJz2"
}

func makeRequest(body []byte, result *http.Response, jwt string, t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error create request %v", err)
	}

	req.Header.Add("Authorization", jwt)
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

func rangeTestTable(tt testTable, jwt string, t *testing.T) {
	for _, v := range tt {
		time.Sleep(v.delay)
		makeRequest(v.in, v.out, jwt, t)
	}
}

func makeRequestWithChecker(body []byte, result *http.Response, jwt string, ch Checker, t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, RPC_ENPOINT, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error create request %v", err)
	}

	req.Header.Add("Authorization", jwt)
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

func rangeTestTableWithChecker(tt testTable, jwt string, ch Checker, t *testing.T) {
	for _, v := range tt {
		time.Sleep(v.delay)
		makeRequestWithChecker(v.in, v.out, jwt, ch, t)
	}
}

func authenticateTestUser(uID, pass string) (string, error) {
	r := strings.NewReader(`{
	  "client_id": "5OCJYuTq5Dog960c3lfVEsBlquDX9Ka2",
	  "username": "` + uID + `",
	  "password": "` + pass + `",
	  "connection": "Username-Password-Authentication",
	  "scope": "openid"
	}`)

	rsp, err := http.DefaultClient.Post("https://kazoup.eu.auth0.com/oauth/ro", "application/json", r)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var a, d authRsp

	if err := json.NewDecoder(rsp.Body).Decode(&a); err != nil {
		return "", err
	}

	data := url.Values{}
	data.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Add("target", "6zIDm8InhbTRp1bL2C4m1TK4Llr4arTy")
	data.Add("scope", "openid")
	data.Add("client_id", "5OCJYuTq5Dog960c3lfVEsBlquDX9Ka2")
	data.Add("api_type", "app")
	data.Add("id_token", a.IDToken)

	rsp2, err := http.DefaultClient.Post("https://kazoup.eu.auth0.com/delegation", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer rsp2.Body.Close()

	if err := json.NewDecoder(rsp2.Body).Decode(&d); err != nil {
		return "", err
	}

	return d.IDToken, nil
}
