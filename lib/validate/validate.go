package validate

import (
	platform_err "github.com/kazoup/platform/lib/errors"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

func Exists(ctx context.Context, params ...string) error {
	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_err.ErrInvalidCtx
	}

	for _, v := range params {
		if len(v) == 0 {
			return errors.BadRequest(srv.Options().Server.String(), "missing parameter")
		}
	}

	return nil
}
