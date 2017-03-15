package subscriber

import (
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

func (a *Announce) Subscriber(ctx context.Context, msg *announce_msg.AnnounceMessage) error {

	log.Println("DATASOURCE ANNOUNCE SUBSCRIBER", msg)

	if globals.HANDLER_DATASOURCE_CREATE == msg.Handler ||
		globals.HANDLER_DATASOURCE_SCAN == msg.Handler {
		log.Println("SCAN NOW")
	}

	return nil
}
