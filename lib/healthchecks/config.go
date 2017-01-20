package healthchecks

import (
	"fmt"
	"github.com/micro/go-os/monitor"
	"net/http"
	"strconv"
)

func RegisterConfigWebHealthChecks(m monitor.Monitor) {
	configWebHealthCheck(m)
}

func configWebHealthCheck(m monitor.Monitor) {
	url := "https://web.kazoup.io:8082/config/health"
	n := "com.kazoup.web.config.health"

	chc := m.NewHealthChecker(
		n,
		"Checking config-web health",
		func() (map[string]string, error) {
			c := &http.Client{}
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				return map[string]string{
					"info": "Error building request",
				}, err
			}
			rsp, err := c.Do(req)
			if err != nil {
				return map[string]string{
					"info": fmt.Sprintf("GET request failed: %s", err),
				}, err
			}

			return map[string]string{
				"info":   "OK",
				"status": strconv.Itoa(rsp.StatusCode),
			}, nil
		},
	)

	if err := m.Register(chc); err != nil {
		fmt.Println("ERROR registering HealthChecker %v", n, err)
	}
}
