package handler

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	proto "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

// Quota struct
type Quota struct {
	Client client.Client
}

// Read quota handler
func (q *Quota) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.quota.Read", err.Error())
	}

	srvs, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.quota.Read", err.Error())
	}

	for _, v := range srvs {
		r, rt, q := quota.GetQuota(v.Name, uID)
		rsp.Quota = append(rsp.Quota, &proto.Quota{
			Name:           v.Name,
			Rate:           r,
			ResetTimestamp: rt,
			Quota:          q,
		})
	}

	log.Println(rsp.Quota)

	return nil
}

// Health
func (q *Quota) Health(ctx context.Context, req *proto.HealthRequest, rsp *proto.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
