package objectstorage

import (
	"bytes"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/objectstorage"
	"io"
	"io/ioutil"
	"time"
)

type MockStorage struct{}

func init() {
	objectstorage.Register(new(MockStorage))
}

func (ms *MockStorage) Init() error {
	return nil
}

func (ms *MockStorage) CreateBucket(bucketName string) error {
	return nil
}

func (ms *MockStorage) DeleteBucket(bucketName string) error {
	return nil
}

func (ms *MockStorage) Upload(r io.ReadCloser, bucketName, key string) error {
	return nil
}

func (ms *MockStorage) Download(bucketName, key string, opts ...string) (io.ReadCloser, error) {
	// platform/auth/web/handler/helper.go
	if key == "test_uuid" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "test_user",
			"exp": time.Date(2050, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		// RPC API ClientID can be found in https://manage.auth0.com/#/clients
		decoded, _ := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
		tokenString, err := token.SignedString(decoded)
		if err != nil {
			return nil, err
		}

		return ioutil.NopCloser(bytes.NewBufferString(tokenString)), nil
	}

	return readCloser{}, nil
}

func (ms *MockStorage) Delete(bucketName, key string) error {
	return nil
}

func (ms *MockStorage) SignedObjectStorageURL(bucketName string, objName string) (string, error) {
	return globals.Mock, nil
}

// helper
type readCloser struct{}

func (rc readCloser) Read(p []byte) (int, error) {
	return 0, nil
}

func (rc readCloser) Close() (err error) {
	return nil
}
