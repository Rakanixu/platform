package tika

import (
	"encoding/json"
	"fmt"
	normalize_text "github.com/kazoup/platform/lib/normalization/text"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

	url := os.Getenv("TIKA_URL")
	if url == "" {
		url = "http://localhost:9998"
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/rmeta", url), rc)
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

	tc[0].XTIKAContent, err = normalize_text.Normalize(tc[0].XTIKAContent)
	if err != nil {
		return nil, err
	}

	return tc[0], nil
}

// ExtractContent receives a io.ReadCloser and returns a Tika interface
func ExtractPlainContent(rc io.ReadCloser) (Tika, error) {
	defer rc.Close()

	url := os.Getenv("TIKA_URL")
	if url == "" {
		url = "http://localhost:9998"
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tika", url), rc)
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

	c, err := normalize_text.NormalizeBytes(b)
	if err != nil {
		return nil, err
	}

	t := &TikaContent{
		XTIKAContent: c,
	}

	return t, nil
}
