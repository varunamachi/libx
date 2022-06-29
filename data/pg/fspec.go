package pg

import (
	"context"
	"fmt"

	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

func getFilterValues(
	gtx context.Context,
	dtype string,
	specs data.FilterSpecList,
	filter *data.Filter) (*data.FilterValues, error) {

	fvals := data.NewFilterValues()
	sel := NewSelectorGenerator().Selector(filter)
	for _, spec := range specs {
		switch spec.Type {
		case data.FtProp:
			fallthrough
		case data.FtArray:
			vals, err := getValues(gtx, dtype, spec, sel)
			if err != nil {
				return nil, err
			}
			fvals.Values[spec.Field] = vals
		case data.FtDateRange:
			dts, err := getDateRangeExtremes(gtx, dtype, spec, sel)
			if err != nil {
				return nil, err
			}
			fvals.Dates[spec.Field] = dts
		case data.FtNumRange:
			rg, err := getRangeExtremes(gtx, dtype, spec, sel)
			if err != nil {
				return nil, err
			}
			fvals.Ranges[spec.Name] = rg

		}
	}

	return nil, nil
}

func getValues(
	gtx context.Context,
	dtype string,
	spec *data.FilterSpec,
	sel Selector) ([]interface{}, error) {

	query := "SELECT DISTINCT %s FROM %s"
	if !sel.IsEmpty() {
		query += " WHERE " + sel.QueryFragment
	}
	query += " ORDER BY %s"
	query = fmt.Sprintf(query, spec.Field, dtype, spec.Field)

	out := make([]interface{}, 0, 100)
	err := Conn().SelectContext(gtx, &out, query, sel.Args...)
	if err != nil {
		return nil, errx.Errf(err,
			"failed to get distinct values for '%s' in '%s'", spec.Field, dtype)
	}
	return out, nil
}

func getDateRangeExtremes(
	gtx context.Context,
	dtype string,
	spec *data.FilterSpec,
	sel Selector) (*data.DateRange, error) {

	query := `SELECT min(%s) as _from, max(%s) as _to FROM %s`
	if !sel.IsEmpty() {
		query += " WHERE " + sel.QueryFragment
	}

	out := data.DateRange{}
	query = fmt.Sprintf(query, spec.Field, spec.Field, dtype)
	err := Conn().SelectContext(gtx, &out, query, sel.Args...)
	if err != nil {
		return nil, errx.Errf(err,
			"failed to get date range for '%s' in '%s'", spec.Field, dtype)
	}
	return &out, nil
}

func getRangeExtremes(
	gtx context.Context,
	dtype string,
	spec *data.FilterSpec,
	sel Selector) (*data.NumberRange, error) {
	query := `SELECT min(%s) as _from, max(%s) as _to FROM %s`
	if !sel.IsEmpty() {
		query += " WHERE " + sel.QueryFragment
	}

	out := data.NumberRange{}
	query = fmt.Sprintf(query, spec.Field, spec.Field, dtype)
	err := Conn().SelectContext(gtx, &out, query, sel.Args...)
	if err != nil {
		return nil, errx.Errf(err,
			"failed to get number range for '%s' in '%s'", spec.Field, dtype)
	}
	return &out, nil
}
