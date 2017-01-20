package healthchecks

import (
	"bytes"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro"
	"github.com/micro/go-os/monitor"
	"net/http"
	"strconv"
)

func RegisterFileHealthChecks(srv micro.Service, m monitor.Monitor) {
	fileSrvHealthCheck(srv, m)
}

func fileSrvHealthCheck(srv micro.Service, m monitor.Monitor) {
	url := "https://web.kazoup.io:8082/rpc"
	body := []byte(`{
		"service":"` + srv.Server().Options().Name + `",
		"method":"File.Health",
		"request":{}
	}`)
	n := fmt.Sprintf("%s.health", srv.Server().Options().Name)

	dshc := m.NewHealthChecker(
		n,
		"Checking file-srv health",
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
		fmt.Println("ERROR registering HealthChecker %v", n, err)
	}
}
