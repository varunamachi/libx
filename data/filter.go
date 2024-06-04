package data

import (
	"fmt"
	"time"
)

type Validatable interface {
	IsValid() bool
}

// Filter - generic filter used to filter data in any mongodb collection
type Filter struct {
	Bools     map[string]any               `json:"bools" db:"bools" bson:"bools"`
	Props     map[string]*Matcher          `json:"props" db:"props" bson:"props"`
	Lists     map[string]*Matcher          `json:"lists" db:"lists" bson:"lists"`
	Searches  map[string]*Matcher          `json:"searches" db:"searches" bson:"searches"`
	Constants map[string]*Matcher          `json:"constants" db:"constants" bson:"constants"`
	Dates     map[string]*DateRangeMatcher `json:"dates" db:"dates" bson:"dates"`
	Ranges    map[string]*RangeMatcher     `json:"range" db:"range" bson:"range"`
}

func NewFilter() *Filter {
	return &Filter{
		Bools:     map[string]any{},
		Props:     map[string]*Matcher{},
		Lists:     map[string]*Matcher{},
		Searches:  map[string]*Matcher{},
		Constants: map[string]*Matcher{},
		Dates:     map[string]*DateRangeMatcher{},
		Ranges:    map[string]*RangeMatcher{},
	}
}

func (f *Filter) Bool(key string, value any) *Filter {
	_, ok := value.(bool)
	if value != nil && !ok {
		panic(fmt.Errorf("invalid bool '%s' given for filter", key))
	}
	f.Bools[key] = value
	return f
}

func (f *Filter) PropIn(key string, values ...any) *Filter {
	return f.matcher("prop", f.Props, false, key, values...)
}

func (f *Filter) PropNotIn(key string, values ...any) *Filter {
	return f.matcher("prop", f.Props, true, key, values...)
}

func (f *Filter) ListIn(key string, values ...any) *Filter {
	return f.matcher("list", f.Lists, false, key, values...)
}

func (f *Filter) ListNotIn(key string, values ...any) *Filter {
	return f.matcher("list", f.Lists, true, key, values...)
}

func (f *Filter) SerachIn(key string, values ...any) *Filter {
	return f.matcher("search", f.Searches, false, key, values...)
}

func (f *Filter) SearchNotIn(key string, values ...any) *Filter {
	return f.matcher("search", f.Searches, true, key, values...)
}

func (f *Filter) ConstIn(key string, values ...any) *Filter {
	return f.matcher("const", f.Constants, false, key, values...)
}

func (f *Filter) ConstNotIn(key string, values ...any) *Filter {
	return f.matcher("const", f.Constants, true, key, values...)
}

func (f *Filter) RangeIn(key string, from, to float64) *Filter {
	rng := RangeMatcher{
		NumberRange: NumberRange{
			From: from,
			To:   to,
		},
		Invert: false,
	}
	if !rng.IsValid() {
		panic(fmt.Errorf("invalid range '%f => %f'", from, to))
	}
	return f
}

func (f *Filter) RangeNotIn(key string, from, to float64) *Filter {
	rng := RangeMatcher{
		NumberRange: NumberRange{
			From: from,
			To:   to,
		},
		Invert: true,
	}
	if !rng.IsValid() {
		panic(fmt.Errorf("invalid range '%f => %f'", from, to))
	}
	return f
}

func (f *Filter) DateRangeIn(key string, from, to time.Time) *Filter {
	rng := DateRangeMatcher{
		DateRange: DateRange{
			From: from,
			To:   to,
		},
		Invert: false,
	}
	if !rng.IsValid() {
		panic(fmt.Errorf("invalid range '%v => %v'", from, to))
	}
	return f
}

func (f *Filter) DateRangeNotIn(key string, from, to time.Time) *Filter {
	rng := DateRangeMatcher{
		DateRange: DateRange{
			From: from,
			To:   to,
		},
		Invert: true,
	}
	if !rng.IsValid() {
		panic(fmt.Errorf("invalid range '%v => %v'", from, to))
	}
	return f
}

func (f *Filter) matcher(
	tp string,
	mp map[string]*Matcher,
	invert bool,
	key string,
	values ...any) *Filter {
	// _, ok := value.(bool)
	// Check value
	if values != nil {
		panic(fmt.Errorf("invalid '%s' item '%s' given for filter", tp, key))
	}
	mp[key] = &Matcher{
		Invert: invert,
		Fields: values,
	}
	return f
}

func IsValid[T Validatable](m map[string]T) bool {
	if len(m) == 0 {
		return false
	}

	for _, val := range m {
		if val.IsValid() {
			return true
		}
	}

	return false
}

func (f *Filter) IsValid(name string, val any) bool {
	if len(f.Bools) != 0 {
		return true
	}
	for _, val := range f.Bools {
		if val == nil { // nil represents tristate
			return true
		}
		if _, ok := val.(bool); ok {
			return true
		}
	}

	return IsValid(f.Props) ||
		IsValid(f.Lists) ||
		IsValid(f.Searches) ||
		IsValid(f.Constants) ||
		IsValid(f.Dates) ||
		IsValid(f.Ranges)

}

type FilterValues struct {
	Values map[string][]any        `json:"values" db:"values" bson:"values"`
	Dates  map[string]*DateRange   `json:"dates" db:"dates" bson:"dates"`
	Ranges map[string]*NumberRange `json:"ranges" db:"ranges" bson:"ranges"`
}

func NewFilterValues() *FilterValues {
	return &FilterValues{
		Values: make(map[string][]any),
		Dates:  make(map[string]*DateRange),
		Ranges: make(map[string]*NumberRange),
	}
}
