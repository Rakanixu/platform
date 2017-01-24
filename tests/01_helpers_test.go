package tests

import "net/http"

const (
	RPC_ENPOINT = "https://web.kazoup.io:8082/rpc"
	USER_ID     = "test@kazoup.com"
	USER_PWD    = "ksu4awemtest"
	STATUS_OK   = 200
)

var c *http.Client

func init() {
	c = &http.Client{}
}
