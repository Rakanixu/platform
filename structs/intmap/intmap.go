package intmap

import (
	"encoding/json"
	"strconv"
)

// Intmap struct
type Intmap map[int]string

// MarshalJSON for Intmap type
func (i Intmap) MarshalJSON() ([]byte, error) {
	x := make(map[string]string)
	for k, v := range i {
		x[strconv.FormatInt(int64(k), 10)] = v

	}
	return json.Marshal(x)

}

// UnmarshalJSON for Intmap type
func (i *Intmap) UnmarshalJSON(b []byte) error {
	x := make(map[string]string)
	if err := json.Unmarshal(b, &x); err != nil {
		return err

	}
	*i = make(Intmap, len(x))
	for k, v := range x {
		if ki, err := strconv.ParseInt(k, 10, 32); err != nil {
			return err

		} else {
			(*i)[int(ki)] = v

		}

	}
	return nil

}
