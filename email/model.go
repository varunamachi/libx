package email

import (
	"github.com/varunamachi/libx/str"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Message struct {
	From       string
	To         []string
	Cc         []string
	Bcc        []string
	Attachment []*mail.File // chnage to custom type if required
	Content    string
}

func (m *Message) SetContent(td *str.TemplateDesc) error {
	c, err := str.SimpleTemplateExpand(td)
	if err != nil {
		return err
	}
	m.Content = c
	return nil
}

type Provider interface {
	Send(desc *Message, html bool) error
}
