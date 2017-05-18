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

func RegisterNotificationSrvHealthChecks(srv micro.Service, m monitor.Monitor) {
	notificationSrvHealthCheck(srv, m)
}

func RegisterNotificationWebHealthChecks(m monitor.Monitor) {
	notificationWebHealthCheck(m)
}

func notificationSrvHealthCheck(srv micro.Service, m monitor.Monitor) {
	url := "https://web.kazoup.io:8082/rpc"
	body := []byte(`{
		"service":"` + srv.Server().Options().Name + `",
		"method":"Service.Health",
		"request":{}
	}`)
	n := fmt.Sprintf("%s.health", srv.Server().Options().Name)

	shc := m.NewHealthChecker(
		n,
		"Checking notification-srv health",
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

	if err := m.Register(shc); err != nil {
		fmt.Printf("ERROR registering HealthChecker %s %s", n, err)
	}
}

func notificationWebHealthCheck(m monitor.Monitor) {
	url := "https://web.kazoup.io:8082/notification/health"
	n := "com.kazoup.web.notification.health"

	whc := m.NewHealthChecker(
		n,
		"Checking notification-web health",
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

	if err := m.Register(whc); err != nil {
		fmt.Printf("ERROR registering HealthChecker %s %s", n, err)
	}
}
