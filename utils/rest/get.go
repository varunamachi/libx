package rest

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

func GetCommonParams(etx echo.Context) (*data.CommonParams, error) {
	pmg := httpx.NewParamGetter(etx)

	page := pmg.QueryInt64Or("page", 0)
	pageSize := pmg.QueryInt64Or("pageSize", 0)
	sort := pmg.QueryStrOr("sort", "")
	sortDesc := pmg.QueryBoolOr("sortDesc", false)

	var filter data.Filter
	pmg.QueryJSON("filter", &filter)

	if pmg.HasError() {
		pmg.WriteDetailedError(os.Stdout)
		return nil, pmg.BadReqError()
	}
	return &data.CommonParams{
		Page:           page,
		PageSize:       pageSize,
		Sort:           sort,
		SortDescending: sortDesc,
		Filter:         &filter,
	}, nil
}

func GetFilter(etx echo.Context) (*data.Filter, error) {
	pmg := httpx.NewParamGetter(etx)
	var filter data.Filter
	pmg.QueryJSON("filter", &filter)

	if pmg.HasError() {
		pmg.WriteDetailedError(os.Stdout)
		return nil, pmg.Error()
	}
	return &filter, nil
}

func Get[T any](etx echo.Context, gdr data.GetterDeleter) ([]T, error) {
	dtype := etx.Param("dtype")
	if dtype == "" {
		return nil, errx.BadReq("data type not given")
	}

	cparams, err := GetCommonParams(etx)
	if err != nil {
		return nil, errx.BadReqX(err,
			"failed to get common parameters to get '%s'", dtype)
	}

	out := make([]T, 0, 100)

	err = gdr.Get(etx.Request().Context(), dtype, cparams, []any{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Count(etx echo.Context, gdr data.GetterDeleter) (int64, error) {
	dtype := etx.Param("dtype")
	if dtype == "" {
		return 0, errx.BadReq("data type not given")
	}

	filter, err := GetFilter(etx)
	if err != nil {
		return 0, errx.BadReqX(err,
			"failed to get filter to count in '%s'", dtype)
	}

	return gdr.Count(etx.Request().Context(), dtype, filter)
}

func GetX[T any](etx echo.Context, tableQuery string) ([]T, error) {

	cparams, err := GetCommonParams(etx)
	if err != nil {
		return nil, errx.BadReqX(err,
			"failed to get common parameters to get data")
	}

	sel := pg.NewSelectorGenerator().SelectorX(cparams)
	query := tableQuery + " WHERE " + sel.QueryFragment

	out := make([]T, 0, 100)

	err = pg.Conn().SelectContext(
		etx.Request().Context(), &out, query, sel.Args...)
	if err != nil {
		return nil, errx.Errf(err, "failed to get data")
	}
	return out, nil
}

func GetOne[T any](etx echo.Context, tableQuery string) (T, error) {
	filter, err := GetFilter(etx)
	var out T
	if err != nil {
		return out, errx.BadReqX(err,
			"failed to get filter to count")
	}

	sel := pg.NewSelectorGenerator().Selector(filter)
	query := tableQuery + " WHERE " + sel.QueryFragment

	err = pg.Conn().GetContext(
		etx.Request().Context(), &out, query, sel.Args...)
	if err != nil {
		return out, errx.Errf(err, "failed to get data")
	}
	return out, nil
}
