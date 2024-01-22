package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/lib/pq"
)

type M map[string]any

func (u M) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *M) Scan(value any) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &u)
}

type Arr interface {
	int64 | float64 | bool | []byte | string | time.Time
}

type Vec[T Arr] []T

func (v Vec[T]) Value() (driver.Value, error) {
	return pq.Array(v).Value()
}

func (v *Vec[T]) Scan(value any) error {
	if value == nil {
		return nil
	}
	return pq.Array((*[]T)(v)).Scan(value)
}

func (v Vec[T]) AsSlice() []T {
	return ([]T)(v)
}

type DbJson[T any] struct {
	val *T
}

func (dbj *DbJson[T]) Scan(value any) error {
	if value == nil {
		dbj.val = nil
		return nil
	}
	dbj.val = new(T)

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, dbj.val)
}

func (dbj *DbJson[T]) Value() (driver.Value, error) {
	return json.Marshal(dbj.val)
}

// DateRange - represents date ranges
type DateRange struct {
	// Name string    `json:"name" bson:"name"`
	From time.Time `json:"from" db:"_from" bson:"from"`
	To   time.Time `json:"to" db:"_to" bson:"to"`
}

// IsValid - returns true if both From and To dates are non-zero
func (r *DateRange) IsValid() bool {
	return !r.From.IsZero() &&
		!r.To.IsZero() &&
		(r.To.Equal(r.From) || r.To.After(r.From))
}

func (r *DateRange) Difference() time.Duration {
	return r.To.Sub(r.From)
}

// NumRange - represents real number ranges
type NumberRange struct {
	From float64 `json:"from" db:"_from" bson:"from"`
	To   float64 `json:"to" db:"_to" bson:"to"`
}

// IsValid - returns true if both From and To dates are non-zero
func (r *NumberRange) IsValid() bool {
	return !math.IsNaN(r.From) &&
		!math.IsInf(r.From, -1) &&
		!math.IsInf(r.From, 1) &&
		!math.IsNaN(r.To) &&
		!math.IsInf(r.To, -1) &&
		!math.IsInf(r.To, 1) &&
		r.To >= r.From
}
