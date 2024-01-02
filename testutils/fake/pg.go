package fake

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

const fakeUserSchema = `
CREATE TABLE IF NOT EXISTS fake_user (
	id					INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name				VARCHAR(60) NOT NULL,	
	first_name			VARCHAR(60) NOT NULL,					
	last_name			VARCHAR(60) NOT NULL,		
	email				VARCHAR(60) NOT NULL,	
	age					INT NOT NULL,
	tags				VARCHAR[] DEFAULT '{}',
	status				VARCHAR(20) DEFAULT 'inactive',
	created				TIMESTAMPZ NOT NULL,
	updated				TIMESTAMPZ NOT NULL
)
`

const fakeItemSchema = `
CREATE TABLE IF NOT EXISTS fake_user (
	id					INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name				VARCHAR(60) NOT NULL,					
	description			VARCHAR(120) NOT NULL,		
	created				TIMESTAMPZ NOT NULL,
	updated				TIMESTAMPZ NOT NULL
)
`

var faker = gofakeit.New(39434)

func FillCmd() *cli.Command {
	return pg.Wrap(&cli.Command{
		Name:        "fill-fake-data",
		Description: "Create fake table and fill fake data",
		Action: func(ctx *cli.Context) error {
			return PgCreateFill(ctx.Context)
		},
	})
}

func PgCreateFill(gtx context.Context) error {
	if err := createSchema(gtx); err != nil {
		return err
	}

	if err := fillUserInfo(gtx); err != nil {
		return err
	}

	if err := fillItemInfo(gtx); err != nil {
		return err
	}

	return nil
}

func createSchema(gtx context.Context) error {
	if _, err := pg.Conn().ExecContext(gtx, fakeUserSchema); err != nil {
		return errx.Errf(err, "failed to create fake user table")
	}
	if _, err := pg.Conn().ExecContext(gtx, fakeItemSchema); err != nil {
		return errx.Errf(err, "failed to create fake item table")
	}
	return nil
}

func fillUserInfo(gtx context.Context) error {

	// tx, err := pg.Conn().Begin()
	// if err != nil {
	// 	return errx.Errf(err, "failed to create DB transaction")
	// }

	// ef := func(err error, fmtStr string, args ...any) error {
	// 	if e := tx.Rollback(); e != nil {
	// 		log.Error().Err(err).
	// 			Msg("transaction rollback failed for fake users table")
	// 	}
	// 	return errx.Errf(err, fmtStr, args...)
	// }

	inserter := squirrel.StatementBuilder.Insert("fake_user").Columns(
		"name",
		"first_name",
		"last_name",
		"email",
		"age",
		"tags",
		"status",
		"created",
		"updated",
	)
	var fakeUser FkUser
	for i := 0; i < 5000; i++ {
		if err := faker.Struct(&fakeUser); err != nil {
			return errx.Errf(err, "failed to create a fake user struct")
		}
		inserter.Values(
			fakeUser.Name,
			fakeUser.FirstName,
			fakeUser.LastName,
			fakeUser.Email,
			fakeUser.Age,
			fakeUser.Tags,
			fakeUser.Status,
			fakeUser.Created,
			fakeUser.Updated,
		)

	}

	query, args, err := inserter.ToSql()
	if err != nil {
		return errx.Errf(err, "failed to create fake user insert query")
	}

	if _, err = pg.Conn().ExecContext(gtx, query, args...); err != nil {
		return errx.Errf(err, "failed to execute fake user insert query")
	}

	return nil
}

func fillItemInfo(gtx context.Context) error {
	inserter := squirrel.StatementBuilder.Insert("fake_user").Columns(
		"name",
		"description",
		"created",
		"updated",
	)
	var fakeItem FkItem
	for i := 0; i < 5000; i++ {
		if err := faker.Struct(&fakeItem); err != nil {
			return errx.Errf(err, "failed to create a fake item struct")
		}
		inserter.Values(
			fakeItem.Name,
			fakeItem.Description,
			fakeItem.Created,
			fakeItem.Updated,
		)

	}

	query, args, err := inserter.ToSql()
	if err != nil {
		return errx.Errf(err, "failed to create fake item insert query")
	}

	if _, err = pg.Conn().ExecContext(gtx, query, args...); err != nil {
		return errx.Errf(err, "failed to execute fake item insert query")
	}

	return nil
}
