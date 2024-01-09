package data

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

func (f *Filter) IsValid(name string, val interface{}) bool {
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
	Values map[string][]interface{} `json:"values" db:"values" bson:"values"`
	Dates  map[string]*DateRange    `json:"dates" db:"dates" bson:"dates"`
	Ranges map[string]*NumberRange  `json:"ranges" db:"ranges" bson:"ranges"`
}

func NewFilterValues() *FilterValues {
	return &FilterValues{
		Values: make(map[string][]any),
		Dates:  make(map[string]*DateRange),
		Ranges: make(map[string]*NumberRange),
	}
}
