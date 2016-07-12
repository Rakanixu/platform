package handler

import "encoding/json"

// ResponseAuth ...
type ResponseAuth struct {
	Token string `json:"token"`
}

// GetResponse ...
func (ra *ResponseAuth) GetResponse() string {
	r := ResponseAuth{
		Token: "hexadecimal.string",
	}

	b, _ := json.Marshal(r)

	return string(b)
}
