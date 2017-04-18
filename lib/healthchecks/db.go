package healthchecks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro"
	"github.com/micro/go-os/monitor"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type statusES struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	RelocatingShards            int     `json:"relocating_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	UnassignedShards            int     `json:"unassigned_shards"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}

func RegisterDBHealthChecks(srv micro.Service, m monitor.Monitor) {
	dbConnectionHealthCheck(srv, m)
	dbSrvHealthCheck(srv, m)
}

func dbConnectionHealthCheck(srv micro.Service, m monitor.Monitor) {
	host := os.Getenv("ELASTICSEARCH_URL")
	if host == "" {
		host = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")
	credentials := ""
	if len(username) > 0 && len(password) > 0 {
		credentials = fmt.Sprintf("%s:%s@", username, password)
	}

	domain := strings.Split(host, "://")
	if len(domain) != 2 {
		return
	}

	url := fmt.Sprintf(
		"%s://%s%s/_cluster/health",
		domain[0],
		credentials,
		domain[1],
	)
	n := fmt.Sprintf("%s.elasticsearch.connection", srv.Server().Options().Name)

	chc := m.NewHealthChecker(
		n,
		"Checking Elastic Search status",
		func() (map[string]string, error) {
			var status statusES

			rsp, err := http.DefaultClient.Get(url)
			if err != nil {
				return map[string]string{
					"info": fmt.Sprintf("GET request failed: %s", url),
				}, err
			}

			if err := json.NewDecoder(rsp.Body).Decode(&status); err != nil {
				return map[string]string{
					"info": fmt.Sprintf("Decoding ES status failed: %s", err),
				}, err
			}

			return map[string]string{
				"info":                             "OK",
				"cluster_name":                     status.ClusterName,
				"status":                           status.Status,
				"timed_out":                        strconv.FormatBool(status.TimedOut),
				"number_of_nodes":                  strconv.Itoa(status.NumberOfNodes),
				"number_of_data_nodes":             strconv.Itoa(status.NumberOfDataNodes),
				"active_primary_shards":            strconv.Itoa(status.ActivePrimaryShards),
				"active_shards":                    strconv.Itoa(status.ActiveShards),
				"relocating_shards":                strconv.Itoa(status.RelocatingShards),
				"initializing_shards":              strconv.Itoa(status.InitializingShards),
				"unassigned_shards":                strconv.Itoa(status.UnassignedShards),
				"delayed_unassigned_shards":        strconv.Itoa(status.DelayedUnassignedShards),
				"number_of_pending_tasks":          strconv.Itoa(status.NumberOfPendingTasks),
				"number_of_in_flight_fetch":        strconv.Itoa(status.NumberOfInFlightFetch),
				"task_max_waiting_in_queue_millis": strconv.Itoa(status.TaskMaxWaitingInQueueMillis),
				"active_shards_percent_as_number":  strconv.FormatFloat(status.ActiveShardsPercentAsNumber, 'f', 6, 64),
			}, nil
		},
	)

	if err := m.Register(chc); err != nil {
		fmt.Printf("ERROR registering HealthChecker %s %s", n, err)
	}
}

func dbSrvHealthCheck(srv micro.Service, m monitor.Monitor) {
	url := "https://web.kazoup.io:8082/rpc"
	body := []byte(`{
		"service":"` + srv.Server().Options().Name + `",
		"method":"DB.Health",
		"request":{}
	}`)
	n := fmt.Sprintf("%s.health", srv.Server().Options().Name)

	dshc := m.NewHealthChecker(
		n,
		"Checking db-srv health",
		func() (map[string]string, error) {
			c := &http.Client{}
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			if err != nil {
				return map[string]string{
					"info": "Error building request",
				}, err
			}
			req.Header.Set("Authorization", globals.SYSTEM_TOKEN)
			req.Header.Set("Content-Type", "application/json")
			rsp, err := c.Do(req)
			if err != nil {
				return map[string]string{
					"info": fmt.Sprintf("POST request with body %s failed: %s", string(body), err),
				}, err
			}

			return map[string]string{
				"info":   "OK",
				"status": strconv.Itoa(rsp.StatusCode),
			}, nil
		},
	)

	if err := m.Register(dshc); err != nil {
		fmt.Printf("ERROR registering HealthChecker %s %s", n, err)
	}
}
