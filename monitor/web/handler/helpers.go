package handler

import (
	"fmt"
	"hash/crc32"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	proto "github.com/micro/go-os/monitor/proto"
)

var (
	colours = []string{"blue", "green", "yellow", "purple", "orange"}
)

type sortedServices struct {
	services []*proto.Service
}

type sortedStatuses struct {
	statuses []*proto.Status
}

func (s sortedServices) Len() int {
	return len(s.services)
}

func (s sortedServices) Less(i, j int) bool {
	return s.services[i].Name < s.services[j].Name
}

func (s sortedServices) Swap(i, j int) {
	s.services[i], s.services[j] = s.services[j], s.services[i]
}

func (s sortedStatuses) Len() int {
	return len(s.statuses)
}

func (s sortedStatuses) Less(i, j int) bool {
	return s.statuses[i].Service.Name < s.statuses[j].Service.Name
}

func (s sortedStatuses) Swap(i, j int) {
	s.statuses[i], s.statuses[j] = s.statuses[j], s.statuses[i]
}
func colour(s string) string {
	return colours[crc32.ChecksumIEEE([]byte(s))%uint32(len(colours))]
}

func distanceOfTime(minutes float64) string {
	switch {
	case minutes < 1:
		return fmt.Sprintf("%d secs", int(minutes*60))
	case minutes < 59:
		return fmt.Sprintf("%d minutes", int(minutes))
	case minutes < 90:
		return "about an hour"
	case minutes < 120:
		return "almost 2 hours"
	case minutes < 1080:
		return fmt.Sprintf("%d hours", int(minutes/60))
	case minutes < 1680:
		return "about a day"
	case minutes < 2160:
		return "more than a day"
	case minutes < 2520:
		return "almost 2 days"
	case minutes < 2880:
		return "about 2 days"
	default:
		return fmt.Sprintf("%d days", int(minutes/1440))
	}

	return ""
}

func timeAgo(t int64) string {
	d := time.Unix(t, 0)
	timeAgo := ""
	startDate := time.Now().Unix()
	deltaMinutes := float64(startDate-d.Unix()) / 60.0
	if deltaMinutes <= 523440 { // less than 363 days
		timeAgo = fmt.Sprintf("%s ago", distanceOfTime(deltaMinutes))
	} else {
		timeAgo = d.Format("2 Jan")
	}

	return timeAgo
}

func hostPath(r *http.Request) string {
	if path := r.Header.Get("X-Micro-Web-Base-Path"); len(path) > 0 {
		return path
	}
	return "/"
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/healthchecks", Healthchecks)
	r.HandleFunc("/healthchecks/{id}", Healthcheck)
	r.HandleFunc("/stats", Stats)
	r.HandleFunc("/stats/{service}", ServiceStats)
	r.HandleFunc("/status", Status)
	r.HandleFunc("/status/{service}", ServiceStatus)
	r.HandleFunc("/services", Index)
	r.HandleFunc("/services/{service}", Service)
	return r
}
