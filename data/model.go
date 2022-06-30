package data

import (
	"context"
)

type M map[string]any

//FilterType - Type of filter item
type FilterType string

//FtProp - filter for a value
const FtProp FilterType = "Prop"

//FtArray - filter for an array
const FtArray FilterType = "Array"

//DateRange - filter for date range
const FtDateRange FilterType = "DateRange"

//FtNumRange - filter for real number range
const FtNumRange FilterType = "NumRange"

//FtBoolean - filter for boolean field
const FtBoolean FilterType = "Boolean"

//FtSearch - filter for search text field
const FtSearch FilterType = "Search"

//FtConstant - constant filter value
const FtConstant FilterType = "Constant"

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

func NewFilterValues() *FilterValues {
	return &FilterValues{
		Values: make(map[string][]interface{}),
		Dates:  make(map[string]*DateRange),
		Ranges: make(map[string]*NumberRange),
	}
}

type CommonParams struct {
	Filter   *Filter
	Page     int64
	PageSize int64
	Sort     string
	SortDesc bool
}

func (qp *CommonParams) Offset() int64 {
	return qp.Page * qp.PageSize
}

func (qp *CommonParams) Limit() int64 {
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
		params *CommonParams,
		out interface{}) error

	FilterValues(
		gtx context.Context,
		dtype string,
		specs FilterSpecList,
		filter *Filter) (*FilterValues, error)
}

type GetterDeleter interface {
	Getter
	Deleter
}
