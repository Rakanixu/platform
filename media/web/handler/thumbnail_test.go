package handler

import (
	"fmt"
	_ "github.com/kazoup/platform/lib/objectstorage/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestThumbnailHandler_ServeHTTP(t *testing.T) {
	index := "index_test"
	file := "file_id_test"
	url := fmt.Sprintf("/media/thumbnail?index=%s&file_id=%s", index, file)
	thumbnailHandler := NewThumbnailHandler()

	r := httptest.NewRecorder()
	handler := http.HandlerFunc(thumbnailHandler.ServeHTTP)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(r, req)

	if http.StatusSeeOther != r.Code {
		t.Errorf("Expected status %v, got %v", http.StatusSeeOther, r.Code)
	}
}
