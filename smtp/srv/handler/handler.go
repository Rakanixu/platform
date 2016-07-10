package handler

import (
	proto "github.com/kazoup/smtp/srv/proto/smtp"
	smtp "github.com/kazoup/smtp/srv/smtp"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// SMTP struct
type SMTP struct{}

// Send srv handler
func (s *SMTP) Send(ctx context.Context, req *proto.SendRequest, rsp *proto.SendResponse) error {
	if len(req.Recipient) <= 0 {
		return errors.BadRequest("go.micro.srv.smtp.SMTP.Send", "Recipient required")
	}

	if err := smtp.Send(req.Recipient, req.Subject, req.Body); err != nil {
		return errors.InternalServerError("go.micro.srv.smtp.SMTP.Send", err.Error())
	}

	return nil
}

// SetSettings srv handler
func (s *SMTP) SetSettings(ctx context.Context, req *proto.SetSettingsRequest, rsp *proto.SetSettingsResponse) error {
	if req.EmailHost == "" {
		return errors.BadRequest("go.micro.srv.smtp.SMTP.SetSettings", "'email_host' required")
	}

	if req.EmailHostUser == "" {
		return errors.BadRequest("go.micro.srv.smtp.SMTP.SetSettings", "'email_host_user' required")
	}

	if req.EmailHostPassword == "" {
		return errors.BadRequest("go.micro.srv.smtp.SMTP.SetSettings", "'email_host_password' required")
	}

	if req.EmailHostPort == "" {
		return errors.BadRequest("go.micro.srv.smtp.SMTP.SetSettings", "'email_host_port' required")
	}

	if req.DefaultFromEmail == "" {
		return errors.BadRequest("go.micro.srv.smtp.SMTP.SetSettings", "'default_from_email' required")
	}

	smtp.EmailHost = req.EmailHost
	smtp.EmailHostUser = req.EmailHostUser
	smtp.EmailHostPassword = req.EmailHostPassword
	smtp.EmailHostPort = req.EmailHostPort
	smtp.DefaultFromEmail = req.DefaultFromEmail

	return nil
}
