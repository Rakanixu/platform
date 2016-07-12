package handler

import (
	"net/http"

	api "github.com/micro/micro/api/proto"
	"golang.org/x/net/context"
)

// Auth ...
type Auth struct{}

/*
Read ...

Expects JSON obj:
{
  "token": string
}
*/
func (a *Auth) Read(ctx context.Context, req *api.Request, res *api.Response) error {
	//TODO: implement

	// MOCK START
	res.StatusCode = http.StatusOK
	ra := ResponseAuth{}
	res.Body = ra.GetResponse()

	// MOCK END
	return nil
}
