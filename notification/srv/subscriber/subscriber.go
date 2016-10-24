package subscriber

import (
	"errors"
	"github.com/kazoup/platform/notification/srv/helpers"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/net/context"
	"log"
)

func Notify(ctx context.Context, nMsg *proto.NotificationMessage) error {
	if len(nMsg.UserId) == 0 {
		return errors.New("ERROR UserId empty")
	}

	log.Println("NOTIFY SUBSCRIPTOR", nMsg.UserId, len(helpers.GetSocketClients().Sockets))

	for _, v := range helpers.GetSocketClients().Sockets {
		log.Println(v.ID)
		if v.ID == nMsg.UserId {
			log.Println("NOTIFY", globals.CODE_REFRESH_DS)
			v.Codes <- globals.CODE_REFRESH_DS
			break
		}
	}

	log.Println("NOTIFY SUBSCRIPTOR, OK")
	return nil
}
