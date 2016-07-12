package handler

import (
	proto "github.com/kazoup/platform/elastic/srv/proto/elastic"
	"github.com/micro/go-micro/errors"
)

// DocRefFieldsExists returns an error if DocRef struct has zero value
func DocRefFieldsExists(dr *proto.DocRef) error {
	if len(dr.Index) <= 0 {
		return errors.BadRequest("go.micro.srv.elastic", "Index required")
	}

	if len(dr.Type) <= 0 {
		return errors.BadRequest("go.micro.srv.elastic", "Type required")
	}

	return nil
}
