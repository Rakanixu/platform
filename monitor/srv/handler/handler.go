package handler

import (
	"github.com/kazoup/platform/monitor/srv/monitor"
	proto "github.com/kazoup/platform/monitor/srv/proto/monitor"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Monitor struct{}

func (m *Monitor) HealthChecks(ctx context.Context, req *proto.HealthChecksRequest, rsp *proto.HealthChecksResponse) error {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	hcs, err := monitor.DefaultMonitor.HealthChecks(req.Service, req.Id, req.Status, int(req.Limit), int(req.Offset))
	if err != nil && err == monitor.ErrNotFound {
		return errors.NotFound("go.micro.srv.monitor.Monitor.HealthCheck", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.monitor.Monitor.HealthCheck", err.Error())
	}

	rsp.Healthchecks = hcs
	return nil
}

func (m *Monitor) Services(ctx context.Context, req *proto.ServicesRequest, rsp *proto.ServicesResponse) error {
	services, err := monitor.DefaultMonitor.Services(req.Service)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.monitor.Monitor.Services", err.Error())
	}
	rsp.Services = services
	return nil
}

func (m *Monitor) Stats(ctx context.Context, req *proto.StatsRequest, rsp *proto.StatsResponse) error {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	stats, err := monitor.DefaultMonitor.Stats(req.Service, req.Id, int(req.Limit), int(req.Offset))
	if err != nil {
		return errors.InternalServerError("go.micro.srv.monitor.Monitor.Stats", err.Error())
	}
	rsp.Stats = stats
	return nil
}

func (m *Monitor) Status(ctx context.Context, req *proto.StatusRequest, rsp *proto.StatusResponse) error {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	statuses, err := monitor.DefaultMonitor.Status(req.Service, req.Id, int(req.Limit), int(req.Offset), req.Verbose)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.monitor.Monitor.Status", err.Error())
	}
	rsp.Statuses = statuses
	return nil
}
