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

type EventUserIdType interface {
	~int | ~int64 | ~string
}

type Event[T EventUserIdType] struct {
	Id        uint64           `json:"id" db:"id" bson:"id"`
	Op        string           `json:"op" db:"op" bson:"op"`
	Type      Type             `json:"type" db:"ev_type" bson:"type"`
	UserId    T                `json:"userId" db:"user_id" bson:"userId"`
	CreatedOn time.Time        `json:"createdOn" db:"created_on" bson:"createdOn"`
	Errors    data.Vec[string] `json:"errors" db:"errors" bson:"errors"`
	Metadata  data.M           `json:"metadata" db:"metadata" bson:"metadata"`
}

type Service[T EventUserIdType] interface {
	AddEvent(gtx context.Context, event *Event[T]) error
}

type Adder[T EventUserIdType] struct {
	event   *Event[T]
	service Service[T]
	gtx     context.Context
}

func NewAdder[T EventUserIdType](
	gtx context.Context,
	service Service[T],
	op string,
	userId T,
	md data.M) *Adder[T] {

	return &Adder[T]{
		event: &Event[T]{
			Op:       op,
			UserId:   userId,
			Type:     None,
			Metadata: md,
		},
		service: service,
		gtx:     gtx,
	}
}

func (adder *Adder[T]) SetType(t Type) *Adder[T] {
	adder.event.Type = t
	return adder
}

func (adder *Adder[T]) SetData(md data.M) *Adder[T] {
	adder.event.Metadata = md
	return adder
}

func (adder *Adder[T]) AddData(name string, md any) *Adder[T] {
	if adder.event.Metadata == nil {
		adder.event.Metadata = data.M{
			name: md,
		}
		return adder
	}
	adder.event.Metadata[name] = md
	return adder
}

func (adder *Adder[T]) SetUser(userId T) *Adder[T] {
	adder.event.UserId = userId
	return adder
}

func (adder *Adder[T]) Commit(err error) error {
	if adder.event.Type == None {
		adder.event.Type = data.Qop(err != nil, Error, Success)
		adder.event.Errors = errx.StackArray(err)
	}
	adder.event.CreatedOn = time.Now()
	if e2 := adder.service.AddEvent(adder.gtx, adder.event); e2 != nil {
		log.Error().Err(e2).Msg("failed to add event to system")
	}
	return err
}

func (adder *Adder[T]) Errf(err error, fmtStr string, args ...any) error {
	err = errx.Errf(err, fmtStr, args...)
	return adder.Commit(err)
}
