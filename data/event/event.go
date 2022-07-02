package event

import (
	"context"
	"time"

	"github.com/varunamachi/libx/data"
)

type Type string

const (
	Success Type = "Success"
	Info    Type = "Info"
	Warning Type = "Warning"
	Error   Type = "Error"
)

type Event struct {
	Name   string
	Type   Type
	UserId string
	Time   time.Time
	Error  []string
	Data   data.M
}

type Service interface {
	AddEvent(gtx context.Context, event *Event) error
}

type Adder struct {
}
