package mg

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func WithMongoFlags(defaultDb string, cmd *cli.Command) *cli.Command {
	cmd.Flags = append(cmd.Flags,
		&cli.StringFlag{
			Name: "mongo-url",
			Usage: "URL of the host running mongodb, this has the highest" +
				" precidence. The default database name is still required " +
				"even when using url",
			EnvVars:  []string{"MONGO_URL"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "mongo-host",
			Value:    "localhost",
			Usage:    "Address of the host running mongodb",
			EnvVars:  []string{"MONGO_HOST"},
			Required: false,
		},
		&cli.IntFlag{
			Name:     "mongo-port",
			Value:    27017,
			Usage:    "Port on which mongodb is listening",
			EnvVars:  []string{"MONGO_PORT"},
			Required: false,
		},
		&cli.StringFlag{
			Name:    "mongo-default-db",
			Value:   defaultDb,
			Usage:   "Name of the database to connect by default",
			EnvVars: []string{"MONGO_DEFAULT_DB"},
		},
		&cli.StringFlag{
			Name:     "mongo-user",
			Usage:    "Mongodb user name",
			EnvVars:  []string{"MONGO_USER"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "mongo-pass",
			Value:    "",
			Usage:    "Mongodb password for connection",
			EnvVars:  []string{"MONGO_PASS"},
			Required: false,
		},
	)

	if cmd.Before == nil {
		cmd.Before = requireMongo
	} else {
		otherBefore := cmd.Before
		cmd.Before = func(ctx *cli.Context) (err error) {
			err = requireMongo(ctx)
			if err == nil {
				err = otherBefore(ctx)
			}
			return err
		}
	}

	return cmd
}

func requireMongo(ctx *cli.Context) error {
	gtx := context.TODO()

	url := ctx.String("mongo-url")
	dbName := ctx.String("mongo-default-db")
	if url != "" {
		if err := Connect(gtx, url, dbName); err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
			return err
		}
	} else {
		err := ConnectWithOpts(
			gtx,
			&ConnOpts{
				Host:     ctx.String("mongo-host"),
				Port:     ctx.Int("mongo-port"),
				User:     ctx.String("mongo-user"),
				DbName:   dbName,
				Password: ctx.String("mongo-pass"),
			})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
			return err
		}
	}
	if err := mongoStore.client.Ping(gtx, nil); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}
	log.Info().Msg("Connected to mongodb")
	return nil
}
