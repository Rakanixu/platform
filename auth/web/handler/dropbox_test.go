package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kazoup/platform/lib/globals"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleDropboxLogin(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test_user",
		"exp": time.Date(2050, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// RPC API ClientID can be found in https://manage.auth0.com/#/clients
	decoded, _ := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
	tokenString, err := token.SignedString(decoded)
	if err != nil {
		t.Fatalf("Error generating valid JWT %v", err)
	}

	testData := []struct {
		method     string
		url        string
		statusCode int
	}{
		{
			method:     http.MethodGet,
			url:        fmt.Sprintf("/dropbox/login?jwt=%s", tokenString),
			statusCode: 307, // Redirect
		},
		{
			method:     http.MethodGet,
			url:        fmt.Sprintf("/dropbox/login?jwt=%s", "invalid"),
			statusCode: 200, // Redirect
		},
		{
			method:     http.MethodPost,
			url:        fmt.Sprintf("/dropbox/login?jwt=%s", "invalid"),
			statusCode: 200, // Redirect
		},
	}

	for _, tt := range testData {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatalf("Error building request: %v", err)
		}

		r := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleDropboxLogin)
		handler.ServeHTTP(r, req)

		if tt.statusCode != r.Code {
			t.Errorf("Expected %v, got %v", tt.statusCode, r.Code)
		}
	}
}
