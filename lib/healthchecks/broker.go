package healthchecks

import (
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-os/monitor"
	"net/http"
	"os"
	"strconv"
)

type statusBroker struct {
	NumSubscriptions int     `json:"num_subscriptions"`
	NumCache         int     `json:"num_cache"`
	NumInserts       int     `json:"num_inserts"`
	NumRemoves       int     `json:"num_removes"`
	NumMatches       int     `json:"num_matches"`
	CacheHitRate     float64 `json:"cache_hit_rate"`
	MaxFanout        int     `json:"max_fanout"`
	AvgFanout        float64 `json:"avg_fanout"`
}

func RegisterBrokerHealthChecks(srv micro.Service, m monitor.Monitor) {
	brokerConnectionHealthCheck(srv, m)
}

func brokerConnectionHealthCheck(srv micro.Service, m monitor.Monitor) {
	// This one is set by kubernetes
	rsh := os.Getenv("REGISTRY_SERVICE_HOST")
	if rsh == "" {
		rsh = "localhost"
	}

	url := fmt.Sprintf("http://%s:8222/subscriptionsz", rsh)
	n := fmt.Sprintf("%s.nats.connection", srv.Server().Options().Name)

	bhc := m.NewHealthChecker(
		n,
		"Checking Broker status",
		func() (map[string]string, error) {
			var status statusBroker

			rsp, err := http.DefaultClient.Get(url)
			if err != nil {
				return map[string]string{
					"info": fmt.Sprintf("GET request failed: %s", url),
				}, err
			}

			if err := json.NewDecoder(rsp.Body).Decode(&status); err != nil {
				return map[string]string{
					"info": fmt.Sprintf("Decoding Broker status failed: %s", err),
				}, err
			}

			return map[string]string{
				"info":              "OK",
				"num_subscriptions": strconv.Itoa(status.NumSubscriptions),
				"num_cache":         strconv.Itoa(status.NumCache),
				"num_inserts":       strconv.Itoa(status.NumInserts),
				"num_removes":       strconv.Itoa(status.NumRemoves),
				"num_matches":       strconv.Itoa(status.NumMatches),
				"cache_hit_rate":    strconv.FormatFloat(status.CacheHitRate, 'f', 6, 64),
				"max_fanout":        strconv.Itoa(status.MaxFanout),
				"avg_fanout":        strconv.FormatFloat(status.AvgFanout, 'f', 6, 64),
			}, nil
		},
	)

	if err := m.Register(bhc); err != nil {
		fmt.Println("ERROR registering HealthChecker %v", n, err)
	}
}
