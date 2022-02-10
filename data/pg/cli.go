package pg

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func WithPostgresFlags(cmd *cli.Command) *cli.Command {
	cmd.Flags = append(cmd.Flags,
		&cli.StringFlag{
			Name: "pg-url",
			Usage: "URL of the host running postgres, this has the highest" +
				" precidence",
			EnvVars:  []string{"PG_URL"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-host",
			Value:    "localhost",
			Usage:    "Address of the host running postgres",
			EnvVars:  []string{"PG_HOST"},
			Required: false,
		},
		&cli.IntFlag{
			Name:     "pg-port",
			Value:    5432,
			Usage:    "Port on which postgres is listening",
			EnvVars:  []string{"PG_PORT"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-db",
			Value:    "",
			Usage:    "Database name",
			EnvVars:  []string{"PG_DB"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-user",
			Usage:    "Postgres user name",
			EnvVars:  []string{"PG_USER"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-pass",
			Value:    "",
			Usage:    "Postgres password for connection",
			EnvVars:  []string{"PG_PASS"},
			Required: false,
		},
	)

	if cmd.Before == nil {
		cmd.Before = requirePostgres
	} else {
		otherBefore := cmd.Before
		cmd.Before = func(ctx *cli.Context) (err error) {
			err = requirePostgres(ctx)
			if err == nil {
				err = otherBefore(ctx)
			}
			return err
		}
	}

	return cmd
}

func requirePostgres(ctx *cli.Context) error {
	url := ctx.String("pg-url")
	if url != "" {
		db, err := Connect(url)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
			return err
		}
		SetDefaultConn(db)
	} else {
		db, err := ConnectWithOpts(&ConnOpts{
			Host:     ctx.String("pg-host"),
			Port:     ctx.Int("pg-port"),
			User:     ctx.String("pg-user"),
			DBName:   ctx.String("pg-db"),
			Password: ctx.String("pg-pass"),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
			return err
		}
		SetDefaultConn(db)
	}
	if err := defDB.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}
	return nil
}
