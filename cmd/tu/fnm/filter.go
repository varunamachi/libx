package fnm

import (
	"context"

	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
)

func GenerateRandomUserFilter(gtx context.Context) (*data.Filter, error) {
	fval, err := pg.GetFilterValues(
		gtx, "fake_user", UserFilterSpec, nil)
	if err != nil {
		return nil, err
	}
	return GetRandomFilter(gtx, UserFilterSpec, fval)
}

func GetRandomFilter(
	gtx context.Context,
	filterSpecs []*data.FilterSpec,
	filterValues *data.FilterValues) (*data.Filter, error) {

	filter := data.NewFilter()

	for _, fspec := range filterSpecs {
		switch fspec.Type {
		case data.FtProp:
			addPropsFilters(filter, fspec, filterValues)
		case data.FtArray:
			addArraysFilters(filter, fspec, filterValues)
		case data.FtDateRange:
			addDateRangeFilters(filter, fspec, filterValues)
		case data.FtNumRange:
			addNumRangeFilters(filter, fspec, filterValues)
		case data.FtConstant:
		case data.FtBoolean:
		case data.FtSearch:
			addSearchFilter(filter, fspec, filterValues)
		}

	}
	return filter, nil
}

func addPropsFilters(
	filter *data.Filter, fspec *data.FilterSpec, fvals *data.FilterValues) {

}
func addArraysFilters(
	filter *data.Filter, fspec *data.FilterSpec, fvals *data.FilterValues) {

}
func addNumRangeFilters(
	filter *data.Filter, fspec *data.FilterSpec, fvals *data.FilterValues) {

}
func addDateRangeFilters(
	filter *data.Filter, fspec *data.FilterSpec, fvals *data.FilterValues) {

}

func addSearchFilter(
	filter *data.Filter, fspec *data.FilterSpec, fvals *data.FilterValues) {

}
