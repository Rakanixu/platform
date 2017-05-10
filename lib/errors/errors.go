package errors

import (
	"fmt"
	micro_errors "github.com/micro/go-micro/errors"
	"reflect"
)

var (
	ErrInvalidCtx = micro_errors.New("Cant get srv from context", "", 500)
)

type PlatformError struct {
	Service string
	Task    string
	Detail  string
	Err     error
}

func NewPlatformError(service, task, detail string, err error) error {
	return &PlatformError{
		Service: service,
		Task:    task,
		Detail:  detail,
		Err:     err,
	}
}

func (e *PlatformError) Error() string {
	return fmt.Sprintf("%s %s %s %s", e.Service, e.Task, e.Detail, e.Err.Error())
}

type DiscoveryError struct {
	Entity interface{}
	Detail string
	Err    error
}

func NewDiscoveryError(entity interface{}, detail string, err error) error {
	return &DiscoveryError{
		Entity: entity,
		Detail: detail,
		Err:    err,
	}
}

func (e *DiscoveryError) Error() string {
	t := reflect.TypeOf(e.Entity)

	return fmt.Sprintf("%s %s %s", t.Name(), e.Detail, e.Err.Error())
}
