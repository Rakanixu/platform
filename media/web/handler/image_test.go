package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/kazoup/platform/lib/db/custom/mock"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestImageHandler_loadDatasources(t *testing.T) {
	imageHandler := NewImageHandler()
	imageHandler.loadDatasources(context.TODO())

	for _, v := range imageHandler.fs {

		mockFs, ok := v.(*fs.MockFs)
		if !ok {
			t.Fatal("Unexpected File System type")
		}
		if globals.Mock != mockFs.Endpoint.Url {
			t.Errorf("Unexpected %v, got %v", globals.Mock, mockFs.Endpoint.Url)
		}
	}
}

func TestImageHandler_getFs(t *testing.T) {
	imageHandler := NewImageHandler()
	imageHandler.loadDatasources(context.TODO())
	fileSys := imageHandler.getFs(&file.KazoupMockFile{
		KazoupFile: file.KazoupFile{
			DatasourceId: globals.Mock,
		},
	})

	mockFs, ok := fileSys.(*fs.MockFs)
	if !ok {
		t.Fatal("Unexpected File System type")
	}
	if globals.Mock != mockFs.Endpoint.Url {
		t.Errorf("Unexpected %v, got %v", globals.Mock, mockFs.Endpoint.Url)
	}
}

func TestImageHandler_ServeHTTP(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test_user",
		"exp": time.Date(2050, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// RPC API ClientID can be found in https://manage.auth0.com/#/clients
	decoded, _ := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
	tokenString, err := token.SignedString(decoded)
	if err != nil {
		t.Fatal(err)
	}

	thumbnailHandler := NewThumbnailHandler()

	url := fmt.Sprintf("/preview?token=%s", tokenString)

	r := httptest.NewRecorder()
	handler := http.HandlerFunc(thumbnailHandler.ServeHTTP)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(r, req)
	if r.Code != http.StatusSeeOther {
		t.Errorf("Expected %v, got %v", http.StatusSeeOther, r.Code)
	}
}
