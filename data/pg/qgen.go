package pg

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

func NewPgGetterDeleter(deleteTaskTableName string) data.GetterDeleter {
	return &PgGetterDeleter{}
}

type PgGetterDeleter struct {
}

func (pgd *PgGetterDeleter) Delete(
	gtx context.Context,
	dataType string,
	keyField string,
	keys ...interface{}) error {

	var buf strings.Builder
	buf.WriteString("DELETE FROM ")
	buf.WriteString(dataType)
	buf.WriteString(" WHERE ")
	buf.WriteString(keyField)
	buf.WriteString("IN (?)")

	query, args, err := sqlx.In(buf.String(), keys)
	if err != nil {
		return errx.Errf(err, "")
	}

	query = Conn().Rebind(query)
	if _, err = defDB.ExecContext(gtx, query, args...); err != nil {
		return errx.Errf(err, "")
	}
	return err
}

func (pgd *PgGetterDeleter) GetOne(
	gtx context.Context,
	dataType string,
	keyField string,
	keys []interface{},
	data interface{}) error {

	var buf strings.Builder
	buf.WriteString("SELECT * FROM ")
	buf.WriteString(dataType)
	buf.WriteString(" WHERE ")
	buf.WriteString(keyField)
	buf.WriteString("IN (?)")

	query, args, err := sqlx.In(buf.String(), keys)
	if err != nil {
		return errx.Errf(err, "failed to get values for")
	}

	query = Conn().Rebind(query)
	if _, err = defDB.ExecContext(gtx, query, args...); err != nil {
		return errx.Errf(err, "")
	}
	return err
}

func (pgd *PgGetterDeleter) Count(
	gtx context.Context,
	dtype string,
	filter *data.Filter) (int64, error) {

	// sel := GenQuery(filter, "SELECT COUNT(*) FROM %s", dtype)
	sel := NewSelectorGenerator().Selector(filter)
	sq := squirrel.StatementBuilder.
		Select("COUNT(*)").
		From(dtype).
		Where(sel.QueryFragment, sel.Args...)

	query, args, err := sq.ToSql()
	if err != nil {
		return 0, errx.Errf(err, "failed to build sql query")
	}

	count := int64(0)
	err = defDB.SelectContext(gtx, &count, query, args...)
	// err := defDB.SelectContext(gtx, &count, sel.QueryFragment, sel.Args...)
	if err != nil {
		return 0, errx.Errf(
			err, "failed to get count for data type '%s'", dtype)
	}
	return count, nil
}

func (pgd *PgGetterDeleter) Get(
	gtx context.Context,
	dtype string,
	params *data.CommonParams,
	out interface{}) error {
	// sel := GenQueryX(params, "SELECT * FROM %s", dtype)
	sel := NewSelectorGenerator().SelectorX(params)
	sq := squirrel.StatementBuilder.
		Select("*").
		From(dtype).
		Where(sel.QueryFragment, sel.Args...)
	query, args, err := sq.ToSql()
	if err != nil {
		return errx.Errf(err, "failed to build sql query")
	}

	err = defDB.SelectContext(gtx, out, query, args...)
	// err := defDB.SelectContext(gtx, out, sel.QueryFragment, sel.Args...)
	if err != nil {
		return errx.Errf(err, "failed to get data for type '%s'", dtype)
	}
	return nil
}

func (pgd *PgGetterDeleter) FilterValues(
	gtx context.Context,
	dtype string,
	specs data.FilterSpecList,
	filter *data.Filter) (*data.FilterValues, error) {
	return getFilterValues(gtx, dtype, specs, filter)
}
