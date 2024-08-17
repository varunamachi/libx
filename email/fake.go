package email

import (
	"errors"
	"sync"

	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/iox"
)

// Stores emails in a map
// Per user 3 maps will store To, CC and BCC mails
// Mails with same name will be overwritten

var (
	ErrNoMailsForUser = errors.New("user has no mails")
)

type userMails struct {
	Username string
	To       map[string]*Message
	Cc       map[string]*Message
	Bcc      map[string]*Message
}

func newUserMails(name string) *userMails {
	return &userMails{
		Username: name,
		To:       map[string]*Message{},
		Cc:       map[string]*Message{},
		Bcc:      map[string]*Message{},
	}
}

func (um *userMails) checkAdd(msg *Message) {
	// Mails will not be duplicated across To, Cc, Bcc
	// We overwrite mails with same id in the same category

}

type FakeEmailProvider struct {
	sync.RWMutex
	mails map[string]*userMails
}

func NewFakeEmailProvider() *FakeEmailProvider {
	return &FakeEmailProvider{
		mails: map[string]*userMails{},
	}
}

func (fp *FakeEmailProvider) Print() {
	iox.PrintJSON(fp.mails)
}

func (fp *FakeEmailProvider) Send(msg *Message, html bool) error {

	get := func(name string) *userMails {
		um := fp.mails[name]
		if um == nil {
			um = newUserMails(name)
			fp.mails[name] = um
		}
		return um
	}

	fp.Lock()
	defer fp.Unlock()
	for _, to := range msg.To {
		get(to).To[msg.Id] = msg
	}
	for _, cc := range msg.Cc {
		get(cc).Cc[msg.Id] = msg
	}
	for _, bcc := range msg.Bcc {
		get(bcc).Bcc[msg.Id] = msg
	}

	return nil

}

func (fp *FakeEmailProvider) Get(user, messageId string) (*Message, error) {
	um, found := fp.mails[user]
	if !found {
		return nil,
			errx.Errf(ErrNoMailsForUser, "user '%s' does not have any mails",
				user)
	}
	msg := um.To[messageId]
	if msg == nil {
		errx.Errf(ErrNoMailsForUser,
			"user '%s' does not have any direct mails", user)

	}
	return msg, nil
}

func (fp *FakeEmailProvider) GetCC(user, messageId string) (*Message, error) {
	um, found := fp.mails[user]
	if !found {
		return nil,
			errx.Errf(ErrNoMailsForUser, "user '%s' does not have any mails")
	}
	msg := um.Cc[messageId]
	if msg == nil {
		errx.Errf(ErrNoMailsForUser,
			"user '%s' does not have any mails in CC", user)

	}
	return msg, nil
}

func (fp *FakeEmailProvider) GetBCC(user, messageId string) (*Message, error) {
	um, found := fp.mails[user]
	if !found {
		return nil,
			errx.Errf(ErrNoMailsForUser, "user '%s' does not have any mails")
	}
	msg := um.Bcc[messageId]
	if msg == nil {
		errx.Errf(ErrNoMailsForUser,
			"user '%s' does not have any mails in BCC", user)

	}
	return msg, nil
}

func (fp *FakeEmailProvider) GetAny(user, messageId string) (*Message, error) {
	um, found := fp.mails[user]
	if !found {
		return nil,
			errx.Errf(ErrNoMailsForUser,
				"user '%s' does not have any mails", user)
	}

	msg := um.To[messageId]
	if msg != nil {
		return msg, nil
	}
	msg = um.Cc[messageId]
	if msg != nil {
		return msg, nil
	}
	msg = um.Bcc[messageId]
	if msg != nil {
		return msg, nil
	}

	return nil, errx.Errf(ErrNoMailsForUser,
		"user '%s' does not have any mails", user)

}
