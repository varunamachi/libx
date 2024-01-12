package fnm

import (
	"context"
	"math/rand"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

func GetDataForRandomFilter(
	gtx context.Context,
	dataType string,
	out []any) error {
	filter, err := GenerateRandomFilter(gtx, dataType)
	if err != nil {
		return err
	}

	cp := data.CommonParams{
		Filter:         filter,
		Page:           0,
		PageSize:       0,
		Sort:           "",
		SortDescending: false,
	}

	if err = pg.NewGetterDeleter().Get(gtx, dataType, &cp, out); err != nil {
		return errx.Errf(
			err,
			"failed to get random values from table '%s'",
			dataType)
	}
	return nil
}

func GenerateRandomFilter(
	gtx context.Context, dataType string) (*data.Filter, error) {
	fval, err := pg.GetFilterValues(gtx, dataType, UserFilterSpec, nil)
	if err != nil {
		return nil, err
	}
	return GetRandomFilter(gtx, dataType, UserFilterSpec, fval)
}

func GetRandomFilter(
	gtx context.Context,
	dataType string,
	filterSpecs []*data.FilterSpec,
	filterValues *data.FilterValues) (*data.Filter, error) {

	filter := data.NewFilter()

	for _, fspec := range filterSpecs {
		switch fspec.Type {
		case data.FtProp:
			fallthrough
		case data.FtArray:
			addPropsAndArrayFilters(filter, fspec, filterValues, 5)
		case data.FtDateRange:
			addDateRangeFilters(filter, fspec, filterValues)
		case data.FtNumRange:
			addNumRangeFilters(filter, fspec, filterValues)
		case data.FtConstant:
		case data.FtBoolean:
		case data.FtSearch:
			err := addSearchFilter(gtx, dataType, filter, fspec, filterValues)
			if err != nil {
				return nil, err
			}
		}
	}
	return filter, nil
}

func addPropsAndArrayFilters(
	filter *data.Filter,
	fspec *data.FilterSpec,
	fvals *data.FilterValues,
	requiredSelections int) {

	vals := fvals.Values[fspec.Field]
	fields := make([]any, requiredSelections)

	for i := 0; i < requiredSelections; i++ {
		idx := rand.Int31n((int32(len(vals))))
		fields[i] = vals[idx]
	}

	filter.Props[fspec.Field] = &data.Matcher{
		Invert: false,
		Fields: fields,
	}
}

func addNumRangeFilters(
	filter *data.Filter,
	fspec *data.FilterSpec,
	fvals *data.FilterValues) {

	nr := fvals.Ranges[fspec.Field]
	mid := int((nr.From + nr.To) / 2)

	start := float64(int32(nr.From) + rand.Int31n(int32(mid/2)))
	end := float64(int32(nr.From) - rand.Int31n(int32(mid/2)))

	filter.Ranges[fspec.Field] = &data.RangeMatcher{
		NumberRange: data.NumberRange{
			From: start,
			To:   end,
		},
		Invert: false,
	}

}

func addDateRangeFilters(
	filter *data.Filter, fspec *data.FilterSpec, fvals *data.FilterValues) {

	fullDr := fvals.Dates[fspec.Field]

	mid := fullDr.Difference() / 2 / 100

	start := fullDr.From.Add(time.Duration(rand.Int31n(int32(mid / 2))))
	end := fullDr.From.Add(-time.Duration(rand.Int31n(int32(mid / 2))))

	filter.Dates[fspec.Field] = &data.DateRangeMatcher{
		DateRange: data.DateRange{
			From: start,
			To:   end,
		},
		Invert: false,
	}
}

func addSearchFilter(
	gtx context.Context,
	table string,
	filter *data.Filter,
	fspec *data.FilterSpec,
	fvals *data.FilterValues) error {

	sq := squirrel.StatementBuilder.
		Select(fspec.Field).
		Distinct().
		From(table).
		Limit(3)
	query, args, err := sq.ToSql()
	if err != nil {
		return errx.Errf(err, "failed to generate query to get distinct "+
			"values for generating search fiter")
	}

	out := make([]interface{}, 0, 100)
	err = pg.Conn().SelectContext(gtx, &out, query, args...)
	if err != nil {
		return errx.Errf(err, "failed to get distinct values for field "+
			"'%s' from table '%s'", fspec.Field, table)
	}

	filter.Searches[fspec.Field] = &data.Matcher{
		Invert: false,
		Fields: out,
	}
	return nil
}
