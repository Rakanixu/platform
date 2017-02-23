package rossete

import (
	"bytes"
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	"net/http"
)

type RosseteEntity struct {
	Type       string `json:"type"`
	Mention    string `json:"mention"`
	Normalized string `json:"normalized"`
	Count      int    `json:"count"`
	EntityID   string `json:"entityId"`
}

type RosseteEntities struct {
	Entities []RosseteEntity `json:"entities"`
}

// Entities queries Rossete API to extract entities from a text
func Entities(text string) (*RosseteEntities, error) {
	// https://developer.rosette.com/api-doc#!/entities/runEntityExtraction
	req, err := http.NewRequest(http.MethodPost, globals.ROSETTE_ENDPOINT, bytes.NewBuffer([]byte(`{"content":"`+text+`"}`)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-RosetteAPI-Key", globals.ROSETTE_API_KEY)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var re *RosseteEntities
	if err := json.NewDecoder(rsp.Body).Decode(&re); err != nil {
		return nil, err
	}

	return re, nil
}
