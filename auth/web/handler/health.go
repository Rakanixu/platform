package handler

import (
	"net/http"
)

// HandleHealthCheck
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Do nothing, will check response status to determine auth-web availability
}
