package countvalue

import (
	"github.com/kazoup/platform/lib/intmap"
)

// CountValue data struct
type CountValue struct {
	Count int           `json:"count"`
	Value intmap.Intmap `json:"value"`
}