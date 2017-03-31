package wrappers

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// NewAfterHandlerWrapper returns a handler quota limit per user wrapper
func NewAfterHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return afterHandlerWrapper(h)
	}
}

// afterHandlerWrapper defines a quota wrapper based on quotaLimit per srv+user_id key
func afterHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if err := fn(ctx, req, rsp); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.New("Cannot get service from context")
		}

		defer func() {
			b, err := json.Marshal(req.Request())
			if err != nil {
				log.Println("ERROR afterHandlerWrapper", err)
			}

			// Publish annuncment after handler was called
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
				globals.AnnounceTopic,
				&announce_msg.AnnounceMessage{
					Handler: fmt.Sprintf("%s.%s", req.Service(), req.Method()),
					Data:    string(b),
				},
			)); err != nil {
				log.Println("ERROR afterHandlerWrapper publishing announcement", err)
			}
		}()

		return nil
	}
}

// NewQuotaSubscriberWrapper returns a subscriber quota limit per user wrapper
func NewAfterSubscriberWrapper() server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return afterSubscriberWrapper(fn)
	}
}

// quotaSubscriberWrapper defines a quota wrapper based on quotaLimit per srv+user_id key
func afterSubscriberWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		if err := fn(ctx, msg); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.New("Cannot get service from context")
		}

		// Announce if topic is not an announcement, we do not want to announce that an announcement was announce..
		if !(msg.Topic() == globals.AnnounceTopic) {
			defer func() {
				b, err := json.Marshal(msg.Message())
				if err != nil {
					log.Println("ERROR afterHandlerWrapper", err)
				}

				// Publish annauncement after subscriber was called
				if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
					globals.AnnounceTopic,
					&announce_msg.AnnounceMessage{
						Handler: msg.Topic(),
						Data:    string(b),
					},
				)); err != nil {
					log.Println("ERROR afterSubscriberWrapper publishing announcement", err)
				}
			}()
		}

		return nil
	}
}
