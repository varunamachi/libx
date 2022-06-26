package pg

import (
	"context"
	"strings"

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
		return errx.Errf(err, "")
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
	// query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s",
	// 	dtype, generateSelector(filter))
	// err = defDB.SelectContext(gtx, &count, query)
	// return count, err
	return 0, nil
}

func (pgd *PgGetterDeleter) Get(
	gtx context.Context,
	dtype string,
	params data.QueryParams,
	out interface{}) error {
	return nil
}

func (pgd *PgGetterDeleter) FilterValues(
	gtx context.Context,
	dtype string,
	field string,
	specs data.FilterSpecList,
	filter *data.Filter) (data.M, error) {
	return nil, nil
}
