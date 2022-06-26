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
		params QueryParams,
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
