# SMTP API [![GoDoc](https://godoc.org/github.com/Rakanixu/smtp/api?status.svg)](https://godoc.org/github.com/Rakanixu/smtp/api)

This is the SMTP API with fqdn go.micro.api.smtp for email delivery.

## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run Service

```
$ go run main.go
```

### Usage

#### Send
```
[POST] http[domain:micro API port]/smtp/send
{
    "recipient":[
        "user@domain.com", 
        "user2@domain.com"
    ], 
    "subject": "Mail subject", 
    "body": "<table style=\"width:100%;\"><tr><td>lets</td><td>see</td></tr><tr><td>the</td><td>markup</td></tr></table>"
}

{}
```

#### Settings
Get resource /smtp/settings definition
```
[OPTIONS] http[domain:micro API port]/smtp/settings
{}

{
  "email_host": "required",
  "email_host_port": "required",
  "email_host_user": "required",
  "email_host_password": "required",
  "default_from_email": "required"
}
```

Create or update SMTP server configuration
```
[PUT] http[domain:micro API port]/smtp/settings
{
    "email_host": "EmailHost",
    "email_host_user": "EmailHostUser",
    "email_host_password": "EmailHostPassword",
    "email_host_port": "EmailHostPort",
    "default_from_email": "DefaultFromEmail"
}

{}
```
