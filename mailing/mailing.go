package mailing

import (
	"bytes"
	"html/template"
	"net/smtp"
)

type TemplateName string

var (
	GMailSMTPServer = "smtp.google.com:587"
)

type Client struct {
	FromEmailAddress string
	SMTPServer       string
	Auth             smtp.Auth

	Templates map[TemplateName]*template.Template
}

func NewMailingClient(smtpAddress string, username string, password string, emailAddress string) *Client {
	auth := smtp.PlainAuth("", username, password, smtpAddress)
	return &Client{
		FromEmailAddress: emailAddress,
		SMTPServer:       smtpAddress,
		Auth:             auth,
		Templates:        map[TemplateName]*template.Template{},
	}
}

func (c *Client) AddTemplate(name TemplateName, path string) error {
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	c.Templates[name] = t
	return nil
}
func (c *Client) SendMail(name TemplateName, replacements any, to ...string) error {
	t := c.Templates[name]

	var body bytes.Buffer

	if err := t.Execute(&body, replacements); err != nil {
		return err
	}
	return smtp.SendMail(c.SMTPServer, c.Auth, c.FromEmailAddress, to, body.Bytes())
}
