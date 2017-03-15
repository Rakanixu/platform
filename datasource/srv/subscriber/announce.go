package subscriber

import (
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type Announce struct {
	Client client.Client
	Broker broker.Broker
}

// Subscriber subscribes to all platform announcment
// Scan have to be trigger as a reaction to creation or scan request
// Scans have to be trigger as a reaction of scan all request
func (a *Announce) Subscriber(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	var es []*proto.Endpoint
	var e *proto.Endpoint
	react := false

	if globals.HANDLER_DATASOURCE_SCAN == msg.Handler {
		var r proto.ScanRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			log.Println(err)
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

		es = append(es, e)

		react = true
	}

	if globals.HANDLER_DATASOURCE_CREATE == msg.Handler {
		var r proto.CreateRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		es = append(es, r.Endpoint)

		react = true
	}

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

		react = true
	}

	if react {
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
