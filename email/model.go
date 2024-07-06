package email

import (
	"errors"

	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
	"github.com/varunamachi/libx/str"
	mail "github.com/xhit/go-simple-mail/v2"
)

var (
	ErrInvalidMailProvider = errors.New("invalid mail provider")
)

type Message struct {
	Id         string
	From       string
	To         []string
	Cc         []string
	Bcc        []string
	Attachment []*mail.File // chnage to custom type if required
	Content    string
	Data       any
}

func (m *Message) SetContent(td *str.TemplateDesc) error {
	c, err := str.SimpleTemplateExpand(td)
	if err != nil {
		return err
	}
	m.Content = c
	m.Data = td.Data
	return nil
}

type Provider interface {
	Send(desc *Message, html bool) error
}

func ProviderFromEnv(envPrefix string) (Provider, error) {
	providerName := rt.EnvString(envPrefix+"_MAIL_PROVIDER", "")
	if providerName == "" {
		return nil, nil
	}

	switch providerName {
	case envPrefix + "_SMTP_MAIL_PROVIDER":
		return NewSMTPProviderFromEnv(envPrefix)
	case envPrefix + "_FAKE_MAIL_PROVIDER":
		return NewFakeEmailProvider(), nil
	case envPrefix + "_SIMPLE_MAIL_SERVICE_CLIENT_PROVIDER":
		return NewSimpleMailSrvClinetFromEnv(envPrefix)
	default:
		return nil, errx.Errf(ErrInvalidMailProvider,
			"unknown mail provider name set as env variable: %s", providerName)

	}
}
