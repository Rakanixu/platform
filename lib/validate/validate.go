package validate

import (
	"github.com/kazoup/platform/lib/errors"
)

func Exists(params ...string) error {
	for _, v := range params {
		if len(v) == 0 {
			return errors.ErrMissingParams
		}
	}

	return nil
}
