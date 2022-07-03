package event

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

type Type string

const (
	Success Type = "Success"
	Info    Type = "Info"
	Warning Type = "Warning"
	Error   Type = "Error"
	None    Type = ""
)

func StrToType(str string) Type {
	switch str {
	case "Success":
		return Success
	case "Info":
		return Info
	case "Warning":
		return Warning
	case "Error":
		return Error
	}
	return None
}

func (t Type) IsValid() bool {
	return t == Success ||
		t == Info ||
		t == Warning ||
		t == Error
}

func (t Type) String() string {
	switch t {
	case Success:
		return "Success"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	}
	return "None"
}

type Event struct {
	Name   string
	UserId string
	Type   Type
	Time   time.Time
	Error  []string
	Data   data.M
}

type Service interface {
	AddEvent(gtx context.Context, event *Event) error
}

type Adder struct {
	event   *Event
	service Service
	gtx     context.Context
}

func NewAdder(
	gtx context.Context,
	service Service,
	name, userId string,
	md data.M) *Adder {

	return &Adder{
		event: &Event{
			Name:   name,
			UserId: userId,
			Type:   None,
			Data:   md,
		},
		service: service,
		gtx:     gtx,
	}
}

func (adder *Adder) SetType(t Type) *Adder {
	adder.event.Type = t
	return adder
}

func (adder *Adder) SetData(md data.M) *Adder {
	adder.event.Data = md
	return adder
}

func (adder *Adder) SetUser(userId string) *Adder {
	adder.event.UserId = userId
	return adder
}

func (adder *Adder) Commit(err error) error {
	if err != nil {
		if adder.event.Type == None {
			adder.event.Type = Error
		}
		adder.event.Error = errx.StackArray(err)
	}
	adder.event.Time = time.Now()
	if e2 := adder.service.AddEvent(adder.gtx, adder.event); e2 != nil {
		log.Error().Err(err).Msg("failed to add event to system")
	}
	return err
}
