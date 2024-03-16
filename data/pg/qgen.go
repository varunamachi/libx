package pg

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

func NewGetterDeleter() data.GetterDeleter {
	return &getterDeleter{}
}

type getterDeleter struct {
}

func (pgd *getterDeleter) Exists(
	gtx context.Context, dtype, keyField string, id any) (bool, error) {

	sq := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("1").
		From(dtype).
		Where(squirrel.Eq{keyField: id}).
		Limit(1)

	sql, args, err := sq.ToSql()
	if err != nil {
		return false, errx.Errf(err,
			"failed to build sql query to check existance of a row")
	}

	query := "SELECT EXISTS(" + sql + ")"

	exists := false
	if err := defDB.GetContext(gtx, &exists, query, args...); err != nil {
		return false, errx.Errf(err,
			"failed to check object existance ('%s' => '%s' == '%s')",
			dtype, keyField, id)
	}
	return exists, nil
}

func (pgd *getterDeleter) Delete(
	gtx context.Context,
	dataType string,
	keyField string,
	keys ...interface{}) error {

	sq := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Delete(dataType).
		Where(squirrel.Eq{keyField: keys})

	query, args, err := sq.ToSql()
	if err != nil {
		return errx.Errf(err, "failed to build sql query")
	}
	if _, err = defDB.ExecContext(gtx, query, args...); err != nil {
		return errx.Errf(err, "failed to delete from %s", dataType)
	}

	return nil
}

func (pgd *getterDeleter) GetOne(
	gtx context.Context,
	dataType string,
	keyField string,
	key interface{},
	data interface{}) error {

	sq := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("*").
		From(dataType).
		Where(squirrel.Eq{keyField: key}).
		Limit(1)

	query, args, err := sq.ToSql()
	if err != nil {
		return errx.Errf(err, "failed to build sql query")
	}
	if err = defDB.GetContext(gtx, data, query, args...); err != nil {
		return errx.Errf(err, "failed to delete from %s", dataType)
	}

	return nil
}

func (pgd *getterDeleter) Count(
	gtx context.Context,
	dtype string,
	filter *data.Filter) (int64, error) {

	sel := NewSelectorGenerator().Selector(filter)
	sq := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("COUNT(*)").
		From(dtype).
		Where(sel.QueryFragment, sel.Args...)

	query, args, err := sq.ToSql()
	if err != nil {
		return 0, errx.Errf(err, "failed to build sql query")
	}

	count := int64(0)
	err = defDB.SelectContext(gtx, &count, query, args...)
	if err != nil {
		return 0, errx.Errf(
			err, "failed to get count for data type '%s'", dtype)
	}
	return count, nil
}

func (pgd *getterDeleter) Get(
	gtx context.Context,
	dtype string,
	params *data.CommonParams,
	out any) error {
	sel := NewSelectorGenerator().SelectorX(params)
	sq := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("*").
		From(dtype).
		Where(sel.QueryFragment, sel.Args...)
	query, args, err := sq.ToSql()

	if err != nil {
		return errx.Errf(err, "failed to build sql query")
	}

	if err = defDB.SelectContext(gtx, out, query, args...); err != nil {
		return errx.Errf(err, "failed to get data for type '%s'", dtype)
	}
	return nil
}

func (pgd *getterDeleter) FilterValues(
	gtx context.Context,
	dtype string,
	specs []*data.FilterSpec,
	filter *data.Filter) (*data.FilterValues, error) {
	return GetFilterValues(gtx, dtype, specs, filter)
}
