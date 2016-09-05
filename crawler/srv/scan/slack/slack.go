package slack

import (
	"github.com/kazoup/platform/crawler/srv/scan"
)

// Slack crawler
type Slack struct {
	Id       int64
	RootPath string
	Index    string
	Running  chan bool
	Scanner  scan.Scanner
}

// Start slack crawler
func (s *Slack) Start(crawls map[int64]scan.Scanner, ds int64) {

}

// Stop slack crawler
func (s *Slack) Stop() {
	s.Running <- false
}

// Info slack crawler
func (s *Slack) Info() (scan.Info, error) {
	return scan.Info{
		Id:          s.Id,
		Type:        "filescanner",
		Description: "Slack scanner",
	}, nil
}
