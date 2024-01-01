package fake

import "time"

type FkUser struct {
	Id        int       `db:"id" fake:"skip"`
	Name      string    `db:"name" fake:"{name}"`
	FirstName string    `db:"firstName" fake:"{firstname}"`
	LastName  string    `db:"lastName" fake:"{lastname}"`
	Email     string    `db:"email" fake:"{email}"`
	Age       int       `db:"age" fake:"{number:1,100}"`
	Tags      []string  `db:"tags" fakesize:"2"`
	Status    string    `db:"status" fake:"{randomstring:[active,inactive]}"`
	Created   time.Time `db:"created"`
	Updated   time.Time `db:"updated"`
}

type FkItem struct {
	Id          int       `db:"id" fake:"skip"`
	Name        string    `db:"name" fake:"{name}"`
	Description string    `db:"description" fake:"{sentence:3}"`
	Created     time.Time `db:"created"`
	Updated     time.Time `db:"updated"`
}
