package categories

import (
	"encoding/json"
)

// Category struct to map file extension into its document type (category)
type Category struct {
	Name       string   `json:"name"`
	Order      int      `json:"order"`
	Extensions []string `json:"extensions"`
	Color      string   `json:"color,omitempty"`
}

var categoryMap []*Category

// SetMap helper
func SetMap(mapping []byte) error {
	if err := json.Unmarshal(mapping, &categoryMap); err != nil {
		return err
	}

	return nil
}

// GetMap helper
func GetMap() []*Category {
	return categoryMap
}

// GetDocType helper
func GetDocType(extension string) string {
	for _, v := range categoryMap {
		if contains(v.Extensions, extension) {
			return v.Name
		}
	}

	return "None"
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
