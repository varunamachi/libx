package email

import (
	"crypto/tls"
	"errors"

	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
	mail "github.com/xhit/go-simple-mail/v2"
)

var (
	ErrInvalidSMTPConfig = errors.New("smtp.invalidConfig")
)

type SmtpConfig struct {
	Host           string
	Port           int
	UserName       string
	Password       string
	SkipEncryption bool
}

type SmtpProvider struct {
	config    SmtpConfig
	tlsConfig *tls.Config
}

func NewSmtpProvider(config SmtpConfig) *SmtpProvider {
	return &SmtpProvider{
		config: config,
	}
}

func (sp *SmtpProvider) WithTLSConfig(cfg *tls.Config) *SmtpProvider {
	sp.tlsConfig = cfg
	return sp
}

func (sp *SmtpProvider) Send(mailDesc *Message, html bool) error {

	m := mail.NewMSG().
		AddTo(mailDesc.To...).
		AddCc(mailDesc.Cc...).
		AddBcc(mailDesc.Bcc...).
		SetFrom(mailDesc.From).
		SetBody(data.Qop(html, mail.TextHTML, mail.TextPlain), mailDesc.Content)

	for _, atc := range mailDesc.Attachment {
		m.Attach(atc)
	}
	if m.Error != nil {
		return errx.Errf(m.Error, "failed to prepare mail for sending")
	}

	client, err := sp.newClient()
	if err != nil {
		return err
	}
	if err := m.Send(client); err != nil {
		return errx.Errf(err, "failed to send email")
	}

	return nil
}

func (sp *SmtpProvider) newClient() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()
	server.Host = sp.config.Host
	server.Port = sp.config.Port
	server.Username = sp.config.UserName
	server.Password = sp.config.Password
	if !sp.config.SkipEncryption {
		server.Encryption = mail.EncryptionSTARTTLS
	}
	if sp.tlsConfig != nil {
		server.TLSConfig = sp.tlsConfig
	}

	client, err := server.Connect()
	if err != nil {
		return nil, errx.Errf(err, "failed to connect to smpt server %s:%d",
			sp.config.Host, sp.config.Port)
	}
	return client, nil
}

func NewSMTPProviderFromEnv(envPrefix string) (Provider, error) {
	host := rt.EnvString(envPrefix+"_SMTP_HOST", "")
	port := rt.EnvInt(envPrefix+"_SMTP_PORT", 0)
	userName := rt.EnvString(envPrefix+"_SMPT_USER", "")
	password := rt.EnvString(envPrefix+"_SMPT_PASSWORD", "")

	if data.OneOf("", host, userName, password) || port == 0 {
		return nil, errx.Errf(ErrInvalidSMTPConfig,
			"unable to get one or more env variables "+
				"required to configure SMTP connection")
	}

	return NewSmtpProvider(SmtpConfig{
		Host:     host,
		Port:     port,
		UserName: userName,
		Password: password,
	}), nil
}
