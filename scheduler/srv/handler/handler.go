package handler

import (
	proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// Scheduler struct
type Scheduler struct {
	Client client.Client
	Crons  []*CronWrapper
}

// CreateScheduledTask srv handler
func (s *Scheduler) CreateScheduledTask(ctx context.Context, req *proto.CreateScheduledTaskRequest, rsp *proto.CreateScheduledTaskResponse) error {
	if len(req.Task.Id) == 0 {
		return errors.BadRequest("go.micro.srv.scheduler", "id required")
	}
	if len(req.Task.Action) == 0 {
		return errors.BadRequest("go.micro.srv.scheduler", "action required")
	}
	if req.Schedule.IntervalSeconds == 0 {
		return errors.BadRequest("go.micro.srv.scheduler", "interval required")
	}

	_, err := s.createTask(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.scheduler", err.Error())
	}

	return nil
}
