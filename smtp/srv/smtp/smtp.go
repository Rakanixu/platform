package smtp

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

// SMTP configuration
var (
	EmailHost         string
	EmailHostPort     string
	EmailHostUser     string
	EmailHostPassword string
	DefaultFromEmail  string
)

// Send mail to all recipients
func Send(recipient []string, subject string, body string) error {
	for _, add := range recipient {
		if err := sendMail(add, subject, body); err != nil {
			return err
		}
	}

	return nil
}

// Opens connection to SMPT server and send single mail
func sendMail(recipient string, subject string, body string) error {
	// Set up authentication information.
	auth := getPlainAuth()

	// Mail address sending the mail
	fromAddress := mail.Address{
		Name:    "",
		Address: DefaultFromEmail,
	}

	// Mail address receiving the mail
	toAddress := mail.Address{
		Name:    "",
		Address: recipient,
	}

	// Mail header
	header := make(map[string]string)
	header["From"] = fromAddress.String()
	header["To"] = toAddress.String()
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Send mail
	err := smtp.SendMail(
		EmailHost+":"+EmailHostPort,
		auth,
		fromAddress.Address,
		[]string{toAddress.Address},
		[]byte(message),
	)

	if err != nil {
		return err
	}

	return nil
}

// Returns an Auth struct
func getPlainAuth() smtp.Auth {
	return smtp.PlainAuth(
		"",
		EmailHostUser,
		EmailHostPassword,
		EmailHost,
	)
}

// Use mail's RFC2047 to encode strings
func encodeRFC2047(s string) string {
	addr := mail.Address{
		Name:    s,
		Address: "",
	}

	return strings.Trim(addr.String(), " <>")
}
