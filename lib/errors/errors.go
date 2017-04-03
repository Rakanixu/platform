package errors

import (
	"github.com/micro/go-micro/errors"
)

var (
	ErrInvalidCtx = errors.New("Cant get srv from context", "", 500)
)
