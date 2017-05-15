package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func Test_PlatformError(t *testing.T) {
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	pErrorType := reflect.TypeOf((*PlatformError)(nil))

	if !pErrorType.Implements(errorType) {
		t.Error("PlatformError does not implement Error")
	}
}

func TestPlatformError_Error(t *testing.T) {
	service := "service"
	task := "task"
	detail := "detail"
	err := errors.New("error")
	e := NewPlatformError(service, task, detail, err)

	if fmt.Sprintf("%s %s %s %s", service, task, detail, err.Error()) != e.Error() {
		t.Errorf("Expected: %v, got: %v", fmt.Sprintf("%s %s %s %s", service, task, detail, err.Error()), e.Error())
	}
}

func Test_DiscoveryError(t *testing.T) {
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	dErrorType := reflect.TypeOf((*DiscoveryError)(nil))

	if !dErrorType.Implements(errorType) {
		t.Error("DiscoveryError does not implement Error")
	}
}

func TestDiscoveryError_Error(t *testing.T) {
	type entity struct {
		Foo string
		Bar interface{}
	}

	ent := entity{
		"string",
		"interface",
	}
	detail := "detail"
	err := errors.New("error")
	e := NewDiscoveryError(ent, detail, err)

	if fmt.Sprintf("entity %s %s", detail, err.Error()) != e.Error() {
		t.Errorf("Expected: %v, got: %v", fmt.Sprintf("entity %s %s", detail, err.Error()), e.Error())
	}
}
