package handler

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	proto "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"sort"
)

// Quota struct
type Quota struct {
	Client client.Client
}

// Search quota handler
func (q *Quota) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.quota.Read", err.Error())
	}

	srvs, err := (*cmd.DefaultOptions().Registry).GetService(req.Srv)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.quota.Read", err.Error())
	}

	for _, v := range srvs {
		l, i, r, rt, q, ok := quota.GetQuota(ctx, v.Name, uID)
		if ok {
			rsp.Quota = &proto.Quota{
				Name:           l,
				Icon:           i,
				Rate:           r,
				ResetTimestamp: rt,
				Quota:          q,
			}
		}
		break
	}

	rsp.TimeLimit = globals.QUOTA_TIME_LIMITER_STRING

	return nil
}

// Search quota handler
func (q *Quota) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.quota.Search", err.Error())
	}

	srvs, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.quota.Search", err.Error())
	}

	var h []*proto.Quota
	for _, v := range srvs {
		l, i, r, rt, q, ok := quota.GetQuota(ctx, v.Name, uID)
		if ok {
			h = append(h, &proto.Quota{
				Name:           l,
				Icon:           i,
				Rate:           r,
				ResetTimestamp: rt,
				Quota:          q,
			})
		}
	}

	// Output will be deterministic
	sort.Sort(sortAlphabetically(h))
	sort.Sort(sortByRate(h))

	rsp.TimeLimit = globals.QUOTA_TIME_LIMITER_STRING
	rsp.Quota = h

	return nil
}

// Health
func (q *Quota) Health(ctx context.Context, req *proto.HealthRequest, rsp *proto.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
