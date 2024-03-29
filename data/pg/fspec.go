package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

func GetFilterValues(
	gtx context.Context,
	dtype string,
	specs []*data.FilterSpec,
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
			fvals.Ranges[spec.Field] = rg

		}
	}

	return fvals, nil
}

func getValues(
	gtx context.Context,
	dtype string,
	spec *data.FilterSpec,
	sel Selector) ([]interface{}, error) {

	sq := squirrel.StatementBuilder.
		Select(spec.Field).
		Distinct().
		From(dtype)
	if !sel.IsEmpty() {
		sq.Where(sel.QueryFragment, sel.Args...)
	}
	query, args, err := sq.ToSql()

	if err != nil {
		return nil, errx.Errf(err, "failed to build sql query")
	}

	// query := "SELECT DISTINCT %s FROM %s"
	// if !sel.IsEmpty() {
	// 	query += " WHERE " + sel.QueryFragment
	// }
	// query += " ORDER BY %s"
	// query = fmt.Sprintf(query, spec.Field, dtype, spec.Field)

	out := make([]interface{}, 0, 100)
	err = Conn().SelectContext(gtx, &out, query, args...)
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

	sq := squirrel.StatementBuilder.Select(
		"min(\""+spec.Field+"\") AS _from",
		"max(\""+spec.Field+"\") AS _to",
	).From(dtype)
	if !sel.IsEmpty() {
		sq.Where(sel.QueryFragment, sel.Args...)
	}

	query, args, err := sq.ToSql()
	if err != nil {
		return nil, errx.Errf(err, "failed to build sql query")
	}

	out := data.DateRange{}
	err = Conn().GetContext(gtx, &out, query, args...)
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

	sq := squirrel.StatementBuilder.Select(
		"min(\""+spec.Field+"\") as _from",
		"max(\""+spec.Field+"\") as _to",
	).From(dtype)
	if !sel.IsEmpty() {
		sq.Where(sel.QueryFragment, sel.Args...)
	}

	query, args, err := sq.ToSql()
	if err != nil {
		return nil, errx.Errf(err, "failed to build sql query")
	}

	out := data.NumberRange{}
	err = Conn().GetContext(gtx, &out, query, args...)
	if err != nil {
		return nil, errx.Errf(err,
			"failed to get number range for '%s' in '%s'", spec.Field, dtype)
	}
	return &out, nil
}
