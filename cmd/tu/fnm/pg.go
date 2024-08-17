package fnm

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/rs/zerolog/log"
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
	created				TIMESTAMPTZ NOT NULL,
	updated				TIMESTAMPTZ NOT NULL
)
`

const fakeItemSchema = `
CREATE TABLE IF NOT EXISTS fake_item (
	id					INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name				VARCHAR(60) NOT NULL,					
	description			VARCHAR(120) NOT NULL,		
	created				TIMESTAMPTZ NOT NULL,
	updated				TIMESTAMPTZ NOT NULL
)
`

var faker = gofakeit.New(39434)

func PgCreateFill(gtx context.Context) error {
	if err := createSchema(gtx); err != nil {
		return errx.Wrap(err)
	}
	log.Info().Msg("schema created successfully")

	if err := clearData(gtx); err != nil {
		return errx.Wrap(err)
	}
	log.Info().Msg("data cleared from all tables")

	if err := fillUserInfo(gtx); err != nil {
		return errx.Wrap(err)
	}
	log.Info().Msg("fake user information created")

	if err := fillItemInfo(gtx); err != nil {
		return errx.Wrap(err)
	}
	log.Info().Msg("fake item data created")

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

func clearData(gtx context.Context) error {
	_, err := pg.Conn().ExecContext(gtx, "DELETE FROM fake_user")
	if err != nil {
		return errx.Errf(err, "failed to delete data from fake_users table")
	}

	_, err = pg.Conn().ExecContext(gtx, "DELETE FROM fake_item")
	if err != nil {
		return errx.Errf(err, "failed to delete data from fake_item table")
	}

	return nil

}

func fillUserInfo(gtx context.Context) error {

	tx, err := pg.Conn().BeginTxx(gtx, nil)
	if err != nil {
		return errx.Errf(err, "failed to create DB transaction")
	}

	ef := func(err error, fmtStr string, args ...any) error {
		if e := tx.Rollback(); e != nil {
			log.Error().Err(err).
				Msg("transaction rollback failed for fake users table")
		}
		return errx.Errf(err, fmtStr, args...)
	}

	var fakeUser FkUser
	for i := 0; i < 5000; i++ {
		if err := faker.Struct(&fakeUser); err != nil {
			return errx.Errf(err, "failed to create a fake user struct")
		}

		inserter := squirrel.StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			Insert("fake_user").
			Columns(
				"name",
				"first_name",
				"last_name",
				"email",
				"age",
				"tags",
				"status",
				"created",
				"updated",
			).
			Values(
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
		query, args, err := inserter.ToSql()
		if err != nil {
			return ef(err, "failed to create fake user insert query")
		}

		if _, err = tx.ExecContext(gtx, query, args...); err != nil {
			return ef(err, "failed to execute fake user insert query")
		}
	}

	if err := tx.Commit(); err != nil {
		return errx.Errf(err, "failed to commit fake-user-add tx")
	}

	return nil
}

func fillItemInfo(gtx context.Context) error {

	tx, err := pg.Conn().BeginTxx(gtx, nil)
	if err != nil {
		return errx.Errf(err, "failed to create DB transaction")
	}

	ef := func(err error, fmtStr string, args ...any) error {
		if e := tx.Rollback(); e != nil {
			log.Error().Err(err).
				Msg("transaction rollback failed for fake users table")
		}
		return errx.Errf(err, fmtStr, args...)
	}

	var fakeItem FkItem
	for i := 0; i < 5000; i++ {
		if err := faker.Struct(&fakeItem); err != nil {
			return ef(err, "failed to create a fake item struct")
		}

		query, args, err := squirrel.
			StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			Insert("fake_item").
			Columns(
				"name",
				"description",
				"created",
				"updated",
			).
			Values(
				fakeItem.Name,
				fakeItem.Description,
				fakeItem.Created,
				fakeItem.Updated,
			).ToSql()
		if err != nil {
			return ef(err, "failed to create fake item insert query")
		}
		if _, err = tx.ExecContext(gtx, query, args...); err != nil {
			return ef(err, "failed to execute fake item insert query")
		}
	}

	if err := tx.Commit(); err != nil {
		return errx.Errf(err, "failed to commit fake-item-add tx")
	}

	return nil
}
