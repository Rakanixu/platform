package healthchecks

import (
	"fmt"
	"github.com/micro/go-os/monitor"
	"net/http"
	"strconv"
)

func RegisterMediaWebHealthChecks(m monitor.Monitor) {
	mediaWebHealthCheck(m)
}

func mediaWebHealthCheck(m monitor.Monitor) {
	url := "https://web.kazoup.io:8082/media/health"
	n := "com.kazoup.web.media.health"

	ahc := m.NewHealthChecker(
		n,
		"Checking media-web health",
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

	if err := m.Register(ahc); err != nil {
		fmt.Println("ERROR registering HealthChecker %v", n, err)
	}
}
