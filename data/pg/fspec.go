package pg

import (
	"context"

	"github.com/varunamachi/libx/data"
)

func getFilterValues(
	gtx context.Context,
	dtype string,
	field string,
	specs data.FilterSpecList,
	filter *data.Filter) (*data.FilterValues, error) {
	return nil, nil
}

func getValues(
	gtx context.Context,
	spec *data.FilterSpec,
	where string) ([]interface{}, error) {
	return nil, nil
}

func getDateRangeExtremes(
	gtx context.Context,
	spec *data.FilterSpec,
	where string) (*data.DateRange, error) {
	return nil, nil
}

func getRangeExtremes(
	gtx context.Context,
	spec *data.FilterSpec,
	where string) (*data.NumberRange, error) {
	return nil, nil
}
