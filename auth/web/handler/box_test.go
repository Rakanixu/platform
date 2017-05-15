package handler

import (
	"encoding/base64"
	//"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/objectstorage/mock"
	//"github.com/kazoup/platform/lib/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var callbackURLs []string

func TestHandleBoxLogin(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test_user",
		"exp": time.Date(2017, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
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
			url:        fmt.Sprintf("/box/login?jwt=%s", tokenString),
			statusCode: 307, // Redirect
		},
		{
			method:     http.MethodGet,
			url:        fmt.Sprintf("/box/login?jwt=%s", "invalid"),
			statusCode: 200, // Redirect
		},
		{
			method:     http.MethodPost,
			url:        fmt.Sprintf("/box/login?jwt=%s", "invalid"),
			statusCode: 200, // Redirect
		},
	}

	for _, tt := range testData {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatalf("Error building request: %v", err)
		}

		r := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleBoxLogin)
		handler.ServeHTTP(r, req)

		if tt.statusCode != r.Code {
			t.Errorf("Expected %v, got %v", tt.statusCode, r.Code)
		}
	}
}

func TestHandleBoxCallback(t *testing.T) {
	/*

		src := []byte("state_info_asdfghjklasdfasdfasdfasdfasdf")
		dst := make([]byte, hex.EncodedLen(len(src)))

		hex.Encode(dst, src)

		state, err := utils.Encrypt(
			[]byte(globals.ENCRYTION_KEY_32),
			dst,
		)
		if err != nil {
			t.Fatal(err)
		}
	*/

	testData := []struct {
		method     string
		url        string
		statusCode int
	}{
	/*		{
			method:     http.MethodPost,
			url:        fmt.Sprintf("/box/callback?state=%s", string(state)),
			statusCode: 200,
		},*/
	/*		{
				method:     http.MethodGet,
				url:        fmt.Sprintf("/box/login?jwt=%s", "invalid"),
				statusCode: 200, // Redirect
			},
			{
				method:     http.MethodPost,
				url:        fmt.Sprintf("/box/login?jwt=%s", "invalid"),
				statusCode: 200, // Redirect
			},*/
	}

	for _, tt := range testData {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatalf("Error building request: %v", err)
		}

		r := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleBoxCallback)
		handler.ServeHTTP(r, req)

		if tt.statusCode != r.Code {
			t.Errorf("Expected %v, got %v", tt.statusCode, r.Code)
		}

		t.Error(r.Body.String())
	}
}
