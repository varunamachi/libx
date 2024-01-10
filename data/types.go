package data

import (
	"math"
	"time"
)

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
