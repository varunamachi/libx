package data

import (
	"context"
)

type M map[string]any

type QueryParams struct {
	Filter   *Filter
	Page     int64
	PageSize int64
	Sort     string
	SortDesc bool
}

func (qp *QueryParams) Offset() int64 {
	return qp.Page * qp.PageSize
}

func (qp *QueryParams) Limit() int64 {
	return qp.PageSize
}

type Deleter interface {
	Delete(
		gtx context.Context,
		dataType string,
		keyField string,
		keys ...interface{}) error
}

type Getter interface {
	Count(
		gtx context.Context,
		dtype string,
		filter *Filter) (int64, error)

	GetOne(
		gtx context.Context,
		dataType string,
		keyField string,
		keys []interface{},
		dataOut interface{}) error

	Get(
		gtx context.Context,
		dtype string,
		params *QueryParams,
		out interface{}) error

	FilterValues(
		gtx context.Context,
		dtype string,
		field string,
		specs FilterSpecList,
		filter *Filter) (values M, err error)
}

type GetterDeleter interface {
	Getter
	Deleter
}

//FilterType - Type of filter item
type FilterType string

//Prop - filter for a value
const Prop FilterType = "Prop"

//Array - filter for an array
const Array FilterType = "arrayArray"

//DateRange - filter for date range
const DtRange FilterType = "DateRange"

//NumRange - filter for real number range
const NumRange FilterType = "NumRange"

//Boolean - filter for boolean field
const Boolean FilterType = "Boolean"

//Search - filter for search text field
const Search FilterType = "Search"

//Constant - constant filter value
const Constant FilterType = "Constant"

//FilterSpec - filter specification
type FilterSpec struct {
	Field string     `json:"field" db:"field" bson:"field"`
	Name  string     `json:"name" db:"name" bson:"name"`
	Type  FilterType `json:"type" db:"type" bson:"type"`
}

//Matcher - matches the given fields. If multiple fileds are given the; the
//joining condition is decided by the MatchStrategy given

//PropMatcher - matches props
type PropMatcher []interface{}

//FilterSpecList - alias for array of filter specs
type FilterSpecList []*FilterSpec

//FilterVal - values for filter along with the count
// type FilterVal struct {
// 	Name  string `json:"name" db:"name" bson:"name"`
// 	Count int    `json:"count" db:"count" bson:"count"`
// }
type Matcher struct {
	Invert bool          `json:"invert" db:"invert" bson:"invert"`
	Fields []interface{} `json:"fields" db:"fields" bson:"fields"`
}

type DateRangeMatcher struct {
	DateRange
	Invert bool `json:"invert" db:"invert" bson:"invert"`
}

type RangeMatcher struct {
	NumberRange
	Invert bool `json:"invert" db:"invert" bson:"invert"`
}

//Filter - generic filter used to filter data in any mongodb collection
type Filter struct {
	Bools     map[string]interface{}       `json:"bools" db:"bools" bson:"bools"`
	Props     map[string]*Matcher          `json:"props" db:"props" bson:"props"`
	Lists     map[string]*Matcher          `json:"lists" db:"lists" bson:"lists"`
	Searches  map[string]*Matcher          `json:"searches" db:"searches" bson:"searches"`
	Constants map[string]*Matcher          `json:"constants" db:"constants" bson:"constants"`
	Dates     map[string]*DateRangeMatcher `json:"dates" db:"dates" bson:"dates"`
	Ranges    map[string]*RangeMatcher     `json:"range" db:"range" bson:"range"`
}

type FilterValues struct {
	Values map[string][]interface{}
	Dates  map[string]*DateRange
	Ranges map[string]*NumberRange
}
