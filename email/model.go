package email

import mail "github.com/xhit/go-simple-mail/v2"

type Content struct {
	Template  string
	Variables map[string]any
}

type Desc struct {
	From       string
	To         []string
	Cc         []string
	Bcc        []string
	Attachment []*mail.File
	Content    string
}

type Provider interface {
	Send(desc *Desc) error
}
