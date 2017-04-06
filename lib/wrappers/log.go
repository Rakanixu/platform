package wrappers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

// LogHandlerWrapper returns a HandlerWrappers that logs all requests
func LogHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return logHandlerWrapper(h)
	}
}

func logHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		var err error
		err = fn(ctx, req, rsp)

		if err != nil {
			log.WithFields(log.Fields{
				"service": req.Service(),
				"handler": req.Method(),
			}).Error(err.Error())
		} else {
			log.WithFields(log.Fields{
				"service": req.Service(),
				"handler": req.Method(),
			}).Info("OK")
		}

		return err
	}
}

// LogHandlerWrapper returns a HandlerWrappers that logs all requests
func LogSubscriberWrapper() server.SubscriberWrapper {
	return func(h server.SubscriberFunc) server.SubscriberFunc {
		return logSubscriberWrapper(h)
	}
}

func logSubscriberWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		var err error
		err = fn(ctx, msg)

		if msg.Topic() != globals.AnnounceTopic {
			uID, err := globals.ParseUserIdFromContext(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"topic": msg.Topic(),
					"task":  "LogSubscriberWrapper",
				}).Error("Unable to retrieve user")
			}

			// Log what happended
			if err != nil {
				log.WithFields(log.Fields{
					"topic": msg.Topic(),
					"user":  uID,
				}).Error(err.Error())
			} else {
				log.WithFields(log.Fields{
					"topic": msg.Topic(),
					"user":  uID,
				}).Info("OK")
			}
		}

		return err
	}
}
