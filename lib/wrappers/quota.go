package wrappers

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	rate "github.com/go-redis/redis_rate"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	timerate "golang.org/x/time/rate"
	"os"
	"time"
)

// NewQuotaHandlerWrapper returns a handler quota limit per user wrapper
func NewQuotaHandlerWrapper(srvName string) server.HandlerWrapper {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "localhost:6379"
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": url,
		},
	})

	limiter := rate.NewLimiter(ring)
	limiter.Fallback = timerate.NewLimiter(timerate.Every(time.Second), 1000)

	return func(h server.HandlerFunc) server.HandlerFunc {
		return quotaHandlerWrapper(h, limiter, srvName)
	}
}

// quotaHandlerWrapper defines a quota wrapper based on quotaLimit per srv+user_id key
func quotaHandlerWrapper(fn server.HandlerFunc, limiter *rate.Limiter, srv string) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if len(md["Authorization"]) == 0 {
			return errors.ErrNoAuthHeader
		}

		// We will read claim to know if public user, or paying or whatever
		token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
			if err != nil {
				return nil, err
			}

			return decoded, nil
		})
		if err != nil {
			return errors.NewPlatformError("", "ParseJWTToken", "", 401, err)
		}

		if token.Claims.(jwt.MapClaims)["roles"] == nil {
			return errors.NewPlatformError("", "ParseJWTToken", "Invalid token", 401, err)
		}

		var quotaLimit int64
		for _, v := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
			switch v.(string) {
			case globals.PRODUCT_TYPE_PERSONAL, globals.PRODUCT_TYPE_TEAM, globals.PRODUCT_TYPE_ENTERPRISE:
				quotaLimit = int64(globals.PRODUCT_QUOTAS.M[v.(string)][srv]["handler"].(int))
			}
		}

		if quotaLimit > 0 {
			_, _, rate, _, quota, ok := quota.Check(ctx, req.Service(), token.Claims.(jwt.MapClaims)["sub"].(string))
			if !ok {
				return errors.ErrQuotaLimitExceeded
			}

			// Quota exceded, respond sync and do not initiate go routines
			if rate-quota > 0 {
				return errors.ErrQuotaLimitExceeded
			}
		}

		// Exec subscriber
		if err := fn(ctx, req, rsp); err != nil {
			return err
		}

		// Update quota once subscriber was succesful
		if quotaLimit > 0 {
			_, _, allowed := limiter.AllowN(fmt.Sprintf("%s-handler-%s", srv, token.Claims.(jwt.MapClaims)["sub"].(string)), quotaLimit, globals.QUOTA_TIME_LIMITER, 1)
			if !allowed {
				return errors.ErrQuotaLimitExceeded
			}
		}

		return nil
	}
}

// NewQuotaSubscriberWrapper returns a subscriber quota limit per user wrapper
func NewQuotaSubscriberWrapper(srvName string) server.SubscriberWrapper {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "localhost:6379"
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": url,
		},
	})
	limiter := rate.NewLimiter(ring)
	limiter.Fallback = timerate.NewLimiter(timerate.Every(time.Second), 1000)

	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return quotaSubscriberWrapper(fn, limiter, srvName)
	}
}

// quotaSubscriberWrapper defines a quota wrapper based on quotaLimit per srv+user_id key
func quotaSubscriberWrapper(fn server.SubscriberFunc, limiter *rate.Limiter, srv string) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		// On announcements just do not apply quota subscriber
		// This is really important as quota counter will be increase when a quated srv listen to global announcements
		if msg.Topic() == globals.AnnounceTopic {
			return fn(ctx, msg)
		}

		md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if len(md["Authorization"]) == 0 {
			return errors.ErrNoAuthHeader
		}

		// We will read claim to know if public user, or paying or whatever
		token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
			if err != nil {
				return nil, err
			}

			return decoded, nil
		})
		if err != nil {
			return errors.NewPlatformError("", "ParseJWTToken", "", 401, err)
		}

		if token.Claims.(jwt.MapClaims)["roles"] == nil {
			return errors.NewPlatformError("", "ParseJWTToken", "Invalid token", 401, err)
		}

		var quotaLimit int64
		for _, v := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
			switch v.(string) {
			case globals.PRODUCT_TYPE_PERSONAL, globals.PRODUCT_TYPE_TEAM, globals.PRODUCT_TYPE_ENTERPRISE:
				quotaLimit = int64(globals.PRODUCT_QUOTAS.M[v.(string)][srv]["subscriber"].(int))
			}
		}

		if quotaLimit > 0 {
			_, _, rate, _, quota, ok := quota.Check(ctx, srv, token.Claims.(jwt.MapClaims)["sub"].(string))
			if !ok {
				return errors.ErrQuotaLimitExceeded
			}

			// Quota exceded, respond sync and do not initiate go routines
			if rate-quota > 0 {
				return errors.ErrQuotaLimitExceeded
			}
		}

		// Exec subscriber
		if err := fn(ctx, msg); err != nil {
			return err
		}

		// Update quota once subscriber task was succesful
		if quotaLimit > 0 {
			_, _, allowed := limiter.AllowN(fmt.Sprintf("%s-subs-%s", srv, token.Claims.(jwt.MapClaims)["sub"].(string)), quotaLimit, globals.QUOTA_TIME_LIMITER, 1)
			if !allowed {
				return errors.ErrQuotaLimitExceeded
			}
		}

		return nil
	}
}
