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
// Scans have to be trigger as a reaction to Datasource creation or scan request
func (a *Announce) Subscriber(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
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

		react = true
	}

	if globals.HANDLER_DATASOURCE_CREATE == msg.Handler {
		var r proto.CreateRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			log.Println(err)
			return err
		}

		e = r.Endpoint

		react = true
	}

	if react {
		if err := a.Client.Publish(ctx, a.Client.NewPublication(
			globals.ScanTopic,
			e,
		)); err != nil {
			return err
		}
	}

	return nil
}
