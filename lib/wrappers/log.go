package wrappers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

// NewLogHandlerWrapper returns a HandlerWrappers that logs all requests
func NewLogHandlerWrapper() server.HandlerWrapper {
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

// NewLogHandlerWrapper returns a HandlerWrappers that logs all requests
func NewLogSubscriberWrapper() server.SubscriberWrapper {
	return func(h server.SubscriberFunc) server.SubscriberFunc {
		return logSubscriberWrapper(h)
	}
}

func logSubscriberWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		var err error
		err = fn(ctx, msg)

		uID, uErr := utils.ParseUserIdFromContext(ctx)
		if uErr != nil {
			log.WithFields(log.Fields{
				"topic": msg.Topic(),
				"task":  "NewLogSubscriberWrapper",
			}).Error("Unable to retrieve user")
		}

		if msg.Topic() != globals.AnnounceTopic {
			// Log what happended
			if err == nil {
				log.WithFields(log.Fields{
					"topic": msg.Topic(),
					"user":  uID,
				}).Info("OK")
			}
		}

		// On Announce Topic, we only check for errors
		// This avoids  to generate to much logs.
		// Succesfull subscribers use to publish a new message specified the task has to be done.
		if err != nil {
			log.WithFields(log.Fields{
				"topic": msg.Topic(),
				"user":  uID,
			}).Error(err.Error())
		}

		return err
	}
}
