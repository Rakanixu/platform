package validate

import (
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

func Exists(ctx context.Context, params ...string) error {
	for _, v := range params {
		if len(v) == 0 {
			return errors.BadRequest("", "missing parameter")
		}
	}

	return nil
}
