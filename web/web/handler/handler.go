package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

type rpcRequest struct {
	Service string
	Method  string
	Address string
	Request interface{}
}

type Handler struct {
	Client client.Client
}

// RPC Handler passes on a JSON or form encoded RPC request to
// a service.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	badRequest := func(description string) {
		e := errors.BadRequest("go.micro.rpc", description)
		w.WriteHeader(400)
		w.Write([]byte(e.Error()))
	}

	var service, method, address string
	var request interface{}

	// response content type
	w.Header().Set("Content-Type", "application/json")

	ct := r.Header.Get("Content-Type")

	// Strip charset from Content-Type (like `application/json; charset=UTF-8`)
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}

	switch ct {
	case "application/json":
		var rpcReq rpcRequest

		d := json.NewDecoder(r.Body)
		d.UseNumber()

		if err := d.Decode(&rpcReq); err != nil {
			badRequest(err.Error())
			return
		}

		service = rpcReq.Service
		method = rpcReq.Method
		address = rpcReq.Address
		request = rpcReq.Request

		// JSON as string
		if req, ok := rpcReq.Request.(string); ok {
			d := json.NewDecoder(strings.NewReader(req))
			d.UseNumber()

			if err := d.Decode(&request); err != nil {
				badRequest("error decoding request string: " + err.Error())
				return
			}
		}
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		method = r.Form.Get("method")

		d := json.NewDecoder(strings.NewReader(r.Form.Get("request")))
		d.UseNumber()

		if err := d.Decode(&request); err != nil {
			badRequest("error decoding request string: " + err.Error())
			return
		}
	}

	if len(service) == 0 {
		badRequest("invalid service")
		return
	}

	if len(method) == 0 {
		badRequest("invalid method")
		return
	}

	// create request/response
	var response json.RawMessage
	var err error

	req := h.Client.NewJsonRequest(service, method, request)

	// create context
	ctx := RequestToContext(r)

	// remote call
	if len(address) > 0 {
		err = h.Client.CallRemote(ctx, address, req, &response)
	} else {
		err = h.Client.Call(ctx, req, &response)
	}
	if err != nil {
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			// assuming it's totally screwed
			ce.Code = 500
			ce.Id = "go.micro.rpc"
			ce.Status = http.StatusText(500)
			ce.Detail = "error during request: " + ce.Detail
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	}

	b, _ := response.MarshalJSON()
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}
