package fnm

import (
	"time"

	"github.com/varunamachi/libx/data"
)

type FkUser struct {
	Id        int              `db:"id" fake:"skip"`
	Name      string           `db:"name" fake:"{name}"`
	FirstName string           `db:"first_name" fake:"{firstname}"`
	LastName  string           `db:"last_name" fake:"{lastname}"`
	Email     string           `db:"email" fake:"{email}"`
	Age       int              `db:"age" fake:"{number:1,100}"`
	Tags      data.Vec[string] `db:"tags" fakesize:"2"`
	Status    string           `db:"status" fake:"{randomstring:[active,inactive]}"`
	Created   time.Time        `db:"created"`
	Updated   time.Time        `db:"updated"`
}

type FkItem struct {
	Id          int       `db:"id" fake:"skip"`
	Name        string    `db:"name" fake:"{name}"`
	Description string    `db:"description" fake:"{sentence:3}"`
	Created     time.Time `db:"created"`
	Updated     time.Time `db:"updated"`
}

var UserFilterSpec = []*data.FilterSpec{
	{Field: "name", Name: "Name", Type: data.FtProp},
	{Field: "first_name", Name: "First Name", Type: data.FtProp},
	{Field: "last_name", Name: "Last Name", Type: data.FtProp},
	{Field: "email", Name: "Email", Type: data.FtProp},
	{Field: "age", Name: "Age", Type: data.FtNumRange},
	{Field: "tags", Name: "Tags", Type: data.FtArray},
	{Field: "status", Name: "Status", Type: data.FtProp},
	{Field: "created", Name: "Created", Type: data.FtDateRange},
	{Field: "updated", Name: "Updated", Type: data.FtDateRange},
}

var ItemFilterSpec = []*data.FilterSpec{
	{Field: "name", Name: "Name", Type: data.FtProp},
	{Field: "description", Name: "Description", Type: data.FtProp},
	{Field: "created", Name: "Created", Type: data.FtProp},
	{Field: "updated", Name: "Updated", Type: data.FtProp},
}
