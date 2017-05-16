package wrappers

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

func NewAuthHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return authHandlerWrapper(h)
	}
}

func authHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		var f error

		md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if len(md["Authorization"]) == 0 {
			return errors.ErrNoAuthHeader
		}

		// Authentication
		if md["Authorization"] != globals.SYSTEM_TOKEN {
			token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.NewPlatformError("", "AuthHandlerWrapper", "Unexpected signing method", nil)
				}

				decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
				if err != nil {
					return nil, err
				}

				return decoded, nil
			})

			if err != nil {
				return errors.NewPlatformError("", "ParseJWTToken", "", err)
			}

			if !token.Valid {
				return errors.NewPlatformError("", "ParseJWTToken", "Invalid token", err)
			}

			ctx = context.WithValue(
				ctx,
				kazoup_context.UserIdCtxKey{},
				kazoup_context.UserIdCtxValue(token.Claims.(jwt.MapClaims)["sub"].(string)),
			)
		}

		f = fn(ctx, req, rsp)

		return f
	}
}

func NewAuthSubscriberWrapper() server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return authSubscriberWrapper(fn)
	}
}

func authSubscriberWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		var f error

		md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if len(md["Authorization"]) == 0 {
			return errors.ErrNoAuthHeader
		}

		// Authentication
		if md["Authorization"] != globals.SYSTEM_TOKEN {
			token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.NewPlatformError("", "AuthHandlerWrapper", "Unexpected signing method", nil)
				}

				decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
				if err != nil {
					return nil, err
				}

				return decoded, nil
			})

			if err != nil {
				return errors.NewPlatformError("", "ParseJWTToken", "", err)
			}

			if !token.Valid {
				return errors.NewPlatformError("", "ParseJWTToken", "Invalid token", err)
			}

			ctx = context.WithValue(
				ctx,
				kazoup_context.UserIdCtxKey{},
				kazoup_context.UserIdCtxValue(token.Claims.(jwt.MapClaims)["sub"].(string)),
			)
		}

		f = fn(ctx, msg)

		return f
	}
}
