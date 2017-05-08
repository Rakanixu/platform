package fs

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
)

func (mfs *MockFs) Authorize() (*proto_datasource.Token, error) {
	return mfs.Endpoint.Token, nil
}

func (mfs *MockFs) GetDatasourceId() string {
	return mfs.Endpoint.Id
}

func (mfs *MockFs) GetThumbnail(id string) (string, error) {
	return globals.Mock, nil
}
