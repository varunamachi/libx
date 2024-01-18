package data

import (
	"context"
)

// FilterType - Type of filter item
type FilterType string

// FtProp - filter for a value
const FtProp FilterType = "Prop"

// FtArray - filter for an array
const FtArray FilterType = "Array"

// DateRange - filter for date range
const FtDateRange FilterType = "DateRange"

// FtNumRange - filter for real number range
const FtNumRange FilterType = "NumRange"

// FtBoolean - filter for boolean field
const FtBoolean FilterType = "Boolean"

// FtSearch - filter for search text field
const FtSearch FilterType = "Search"

// FtConstant - constant filter value
const FtConstant FilterType = "Constant"

// FilterSpec - filter specification
type FilterSpec struct {
	Field string     `json:"field" db:"field" bson:"field"`
	Name  string     `json:"name" db:"name" bson:"name"`
	Type  FilterType `json:"type" db:"type" bson:"type"`
}

//Matcher - matches the given fields. If multiple fileds are given the; the
//joining condition is decided by the MatchStrategy given

// PropMatcher - matches props
type PropMatcher []any

// FilterVal - values for filter along with the count
//
//	type FilterVal struct {
//		Name  string `json:"name" db:"name" bson:"name"`
//		Count int    `json:"count" db:"count" bson:"count"`
//	}
type Matcher struct {
	Invert bool  `json:"invert" db:"invert" bson:"invert"`
	Fields []any `json:"fields" db:"fields" bson:"fields"`
}

func (m *Matcher) IsValid() bool {
	for _, val := range m.Fields {
		if val != nil {
			return true
		}
	}

	return false
}

type DateRangeMatcher struct {
	DateRange
	Invert bool `json:"invert" db:"invert" bson:"invert"`
}

func (dr *DateRangeMatcher) IsValid() bool {
	return dr.DateRange.IsValid()
}

type RangeMatcher struct {
	NumberRange
	Invert bool `json:"invert" db:"invert" bson:"invert"`
}

func (r *RangeMatcher) IsValid() bool {
	return r.NumberRange.IsValid()
}

type CommonParams struct {
	Filter         *Filter `json:"filter" db:"filter" bson:"filter"`
	Page           int64   `json:"page" db:"page" bson:"page"`
	PageSize       int64   `json:"pageSize" db:"page_size" bson:"pageSize"`
	Sort           string  `json:"sort" db:"sort" bson:"sort"`
	SortDescending bool    `json:"sortDescending" db:"sort_desc" bson:"sortDescending"`
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
		keys ...any) error
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
		key any,
		dataOut any) error

	Get(
		gtx context.Context,
		dtype string,
		params *CommonParams,
		out any) error

	FilterValues(
		gtx context.Context,
		dtype string,
		specs []*FilterSpec,
		filter *Filter) (*FilterValues, error)
}

type GetterDeleter interface {
	Getter
	Deleter
}
