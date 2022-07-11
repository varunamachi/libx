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
	Op        string    `json:"op" db:"op" bson:"op"`
	UserId    string    `json:"userId" db:"user_id" bson:"userId"`
	Type      Type      `json:"type" db:"ev_type" bson:"type"`
	CreatedOn time.Time `json:"createdOn" db:"creadted_on" bson:"createdOn"`
	Error     []string  `json:"error" db:"error" bson:"error"`
	Metadata  data.M    `json:"metadata" db:"metadata" bson:"metadata"`
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
	op, userId string,
	md data.M) *Adder {

	return &Adder{
		event: &Event{
			Op:       op,
			UserId:   userId,
			Type:     None,
			Metadata: md,
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
	adder.event.Metadata = md
	return adder
}

func (adder *Adder) AddData(name string, md any) *Adder {
	if adder.event.Metadata == nil {
		adder.event.Metadata = data.M{
			name: md,
		}
		return adder
	}
	adder.event.Metadata[name] = md
	return adder
}

func (adder *Adder) SetUser(userId string) *Adder {
	adder.event.UserId = userId
	return adder
}

func (adder *Adder) Commit(err error) error {
	if adder.event.Type == None {
		adder.event.Type = data.Qop(err != nil, Error, Success)
		adder.event.Error = errx.StackArray(err)
	}
	adder.event.CreatedOn = time.Now()
	if e2 := adder.service.AddEvent(adder.gtx, adder.event); e2 != nil {
		log.Error().Err(e2).Msg("failed to add event to system")
	}
	return err
}
