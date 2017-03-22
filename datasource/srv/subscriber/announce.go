package subscriber

import (
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	deletebucket_msg "github.com/kazoup/platform/lib/protomsg/deletebucket"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type AnnounceDatasource struct {
	Client client.Client
	Broker broker.Broker
}

// OnDatasourceCreate
func (a *AnnounceDatasource) OnDatasourceCreate(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// Trigger scan on datasource creation
	if globals.HANDLER_DATASOURCE_CREATE == msg.Handler {
		var r proto.CreateRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		if err := a.Client.Publish(ctx, a.Client.NewPublication(
			globals.ScanTopic,
			r.Endpoint,
		)); err != nil {
			return err
		}
	}

	return nil
}

// OnDatasourceDelete
func (a *AnnounceDatasource) OnDatasourceDelete(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// Trigger bucket deletion for datasource
	if globals.HANDLER_DATASOURCE_DELETE == msg.Handler {
		var r proto.DeleteRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		if err := a.Client.Publish(ctx, a.Client.NewPublication(
			globals.DeleteBucketTopic,
			&deletebucket_msg.DeleteBucketMsg{
				Index: r.Index,
			},
		)); err != nil {
			return err
		}
	}

	return nil
}

// OnDatasourceScan, trigger scan
func (a *AnnounceDatasource) OnDatasourceScan(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	var e *proto.Endpoint

	if globals.HANDLER_DATASOURCE_SCAN == msg.Handler {
		var r proto.ScanRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		rr, err := db_helper.ReadFromDB(a.Client, ctx, &db_proto.ReadRequest{
			Index: globals.IndexDatasources,
			Type:  globals.TypeDatasource,
			Id:    r.Id,
		})
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(rr.Result), &e); err != nil {
			return err
		}

		if err := a.Client.Publish(ctx, a.Client.NewPublication(
			globals.ScanTopic,
			e,
		)); err != nil {
			return err
		}
	}

	return nil
}

// OnDatasourceScanAll trigger scans
func (a *AnnounceDatasource) OnDatasourceScanAll(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	var es []*proto.Endpoint
	var e *proto.Endpoint

	if globals.HANDLER_DATASOURCE_SCANALL == msg.Handler {
		var r proto.ScanAllRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		if len(r.DatasourcesId) > 0 {
			for _, v := range r.DatasourcesId {
				rr, err := db_helper.ReadFromDB(a.Client, ctx, &db_proto.ReadRequest{
					Index: globals.IndexDatasources,
					Type:  globals.TypeDatasource,
					Id:    v,
				})
				if err != nil {
					return err
				}

				if err := json.Unmarshal([]byte(rr.Result), &e); err != nil {
					return err
				}

				es = append(es, e)
			}

		} else {
			srvRes, err := db_helper.SearchFromDB(a.Client, ctx, &db_proto.SearchRequest{
				Index: globals.IndexDatasources,
				Type:  globals.TypeDatasource,
				From:  0,
				Size:  9999,
			})
			if err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(srvRes.Result), &es); err != nil {
				return err
			}
		}

		for _, v := range es {
			if err := a.Client.Publish(ctx, a.Client.NewPublication(
				globals.ScanTopic,
				v,
			)); err != nil {
				return err
			}
		}
	}

	return nil
}
