package handler

import (
	"github.com/getlantern/context"
	proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
)

type Scheduler struct {
	Client client.Client
}

func (s *Scheduler) CreateScheduledTask(ctx context.Context, req *proto.CreateScheduledTaskRequest, rsp *proto.CreateScheduledTaskResponse) error {

	return nil
}
