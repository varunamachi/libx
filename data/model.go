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

type DataDeleter interface {
	Name() string
	Delete(
		gtx context.Context,
		dataType string,
		keyField string,
		key interface{}) error
}

type DataRetriever interface {
	Name() string
	Count(
		gtx context.Context,
		dtype string,
		filter *Filter) (int64, error)
	RetrieveOne(
		gtx context.Context,
		dataType string,
		keyField string,
		key interface{},
		data interface{}) error
	Retrieve(
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
