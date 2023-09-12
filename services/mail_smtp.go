package services

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"

	"git.softndit.com/collector/backend/config"

	"github.com/BurntSushi/toml"
)

// SMTPMailConfig TBD
type SMTPMailConfig struct {
	Password   string
	Username   string
	ServerName string
	Port       int
}

func (c *SMTPMailConfig) parse(tomlCfg *toml.Primitive) error {
	return toml.PrimitiveDecode(*tomlCfg, c)
}

// ParseServerConfig TBD
func (c *SMTPMailConfig) ParseServerConfig(mc *config.MailClientConfig) error {
	c.Password = mc.Password
	c.Username = mc.Username
	c.ServerName = mc.ServerName
	c.Port = mc.Port

	return nil
}

// SMTPEmail TBD
type SMTPEmail struct {
	auth smtp.Auth
	addr string
	port int
}

// NewSMTPEmail TBD
func NewSMTPEmail(config *SMTPMailConfig) (*SMTPEmail, error) {
	sender := &SMTPEmail{}
	if err := sender.configure(config); err != nil {
		return nil, err
	}
	return sender, nil
}

func (s *SMTPEmail) configure(config *SMTPMailConfig) error {
	s.auth = smtp.PlainAuth(
		"",
		config.Username,
		config.Password,
		config.ServerName)

	s.addr = config.ServerName
	s.port = config.Port

	return nil
}

// Send TBD
func (s *SMTPEmail) Send(mail *Mail) error {
	header := make(map[string]string)
	header["From"] = mail.From
	header["To"] = strings.Join(mail.To, ",")
	header["Subject"] = mail.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(mail.Body))

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.addr, s.port),
		s.auth,
		mail.From,
		mail.To,
		[]byte(message),
	)

	return err
}

var _ MailSender = (*SMTPEmail)(nil)
