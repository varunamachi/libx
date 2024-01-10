package fnm

import (
	"context"
	"math/rand"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
)

func GenerateRandomUserFilter(gtx context.Context) (*data.Filter, error) {
	fval, err := pg.GetFilterValues(
		gtx, "fake_user", UserFilterSpec, nil)
	if err != nil {
		return nil, err
	}
	return GetRandomFilter(gtx, "fake-user", UserFilterSpec, fval)
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
			addSearchFilter(gtx, dataType, filter, fspec, filterValues)
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

	mid := fullDr.Difference() / 2

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
	fvals *data.FilterValues) {

	sq := squirrel.StatementBuilder.
		Select(fspec.Field).
		Distinct().
		From(table).
		Limit(3)
	query, args, err := sq.ToSql()
	if err != nil {
		log.Error().Err(err).Msg(
			"failed to generate query to get distinct values " +
				"for generating search fiter")
	}

	out := make([]interface{}, 0, 100)
	err = pg.Conn().SelectContext(gtx, &out, query, args...)
	if err != nil {
		log.Error().
			Err(err).
			Str("field", fspec.Field).
			Str("dataType", table).
			Msg("failed to get distinct values")
	}

	filter.Searches[fspec.Field] = &data.Matcher{
		Invert: false,
		Fields: out,
	}

}
