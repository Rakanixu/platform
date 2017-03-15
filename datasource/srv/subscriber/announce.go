package subscriber

import (
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

func (a *Announce) Subscriber(ctx context.Context, msg *announce_msg.AnnounceMessage) error {

	log.Println("DATASOURCE ANNOUNCE SUBSCRIBER", msg)

	// On Hanlers Datasource.Create, Datasource.Scan publish scanTopic

	return nil
}
