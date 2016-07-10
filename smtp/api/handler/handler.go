package handler

import (
	"encoding/json"
	"net/http"

	smtp "github.com/kazoup/smtp/srv/proto/smtp"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
	"golang.org/x/net/context"
)

// Smtp struct
type Smtp struct{}

// Send (mail) API handler
func (s *Smtp) Send(ctx context.Context, req *api.Request, rsp *api.Response) error {
	if len(req.Body) <= 0 {
		return errors.BadRequest("go.micro.api.smtp", "Required fields")
	}

	var input map[string]interface{}
	json.Unmarshal([]byte(req.Body), &input)

	var recipient = make([]string, len(input["recipient"].([]interface{})))
	for k, v := range input["recipient"].([]interface{}) {
		recipient[k] = v.(string)
	}

	sendReq := client.NewRequest(
		"go.micro.srv.smtp",
		"SMTP.Send",
		&smtp.SendRequest{
			Recipient: recipient,
			Subject:   input["subject"].(string),
			Body:      input["body"].(string),
		},
	)
	sendRsp := &smtp.SendResponse{}

	if err := client.Call(ctx, sendReq, sendRsp); err != nil {
		return errors.InternalServerError("go.micro.api.smtp", err.Error())
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}

// Settings mail server API handler
func (s *Smtp) Settings(ctx context.Context, req *api.Request, rsp *api.Response) error {
	// TODO: waiting for pull request acceptance, won't work until merged
	// https://github.com/micro/micro/pull/62
	if req.Method == "OPTIONS" {
		options, _ := json.Marshal(&smtp.SetSettingsRequest{
			EmailHost:         "required",
			EmailHostPort:     "required",
			EmailHostUser:     "required",
			EmailHostPassword: "required",
			DefaultFromEmail:  "required",
		})
		rsp.Body = string(options)
	}

	if req.Method == "PUT" {
		var input map[string]interface{}
		json.Unmarshal([]byte(req.Body), &input)

		if input["email_host"] == nil {
			return errors.BadRequest("go.micro.api.smtp", "email_host required")
		}

		if input["email_host_port"] == nil {
			return errors.BadRequest("go.micro.api.smtp", "email_host_port required")
		}

		if input["email_host_user"] == nil {
			return errors.BadRequest("go.micro.api.smtp", "email_host_user required")
		}

		if input["email_host_password"] == nil {
			return errors.BadRequest("go.micro.api.smtp", "email_host_password required")
		}

		if input["default_from_email"] == nil {
			return errors.BadRequest("go.micro.api.smtp", "default_from_email required")
		}

		srvReq := client.NewRequest(
			"go.micro.srv.smtp",
			"SMTP.SetSettings",
			&smtp.SetSettingsRequest{
				EmailHost:         input["email_host"].(string),
				EmailHostPort:     input["email_host_port"].(string),
				EmailHostUser:     input["email_host_user"].(string),
				EmailHostPassword: input["email_host_password"].(string),
				DefaultFromEmail:  input["default_from_email"].(string),
			},
		)
		srvRsp := &smtp.SetSettingsResponse{}

		if err := client.Call(ctx, srvReq, srvRsp); err != nil {
			return errors.InternalServerError("go.micro.api.smtp", err.Error())
		}

		rsp.Body = `{}`
	}

	rsp.StatusCode = http.StatusOK

	return nil
}
