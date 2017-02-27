package tika

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

type Tika interface {
	Content() string
}

type TikaContent struct {
	Author       string `json:"Author"`
	LastAuthor   string `json:"Last-Author"`
	WordCount    string `json:"Word-Count"`
	XTIKAContent string `json:"X-TIKA:content"`
	CpRevision   string `json:"cp:revision"`
}

func (tc *TikaContent) Content() string {
	return tc.XTIKAContent
}

// ExtractContent receives a io.ReadCloser and returns a Tika interface
func ExtractContent(rc io.ReadCloser) (Tika, error) {
	defer rc.Close()

	req, err := http.NewRequest(http.MethodPut, "http://tika:9998/rmeta", rc)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var tc []*TikaContent

	if err := json.NewDecoder(rsp.Body).Decode(&tc); err != nil {
		return nil, err
	}

	if len(tc) == 0 {
		return nil, errors.New("No result")
	}

	return tc[0], nil
}

// ExtractContent receives a io.ReadCloser and returns a Tika interface
func ExtractPlainContent(rc io.ReadCloser) (Tika, error) {
	defer rc.Close()

	req, err := http.NewRequest(http.MethodPut, "http://tika:9998/tika", rc)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "text/plain")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	t := &TikaContent{
		XTIKAContent: string(b),
	}

	return t, nil
}
