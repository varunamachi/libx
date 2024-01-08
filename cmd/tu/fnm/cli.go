package fnm

import (
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/iox"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		fillFakeDataCmd(),
		printFValsForFakeData(),
	}
}

func fillFakeDataCmd() *cli.Command {
	return pg.Wrap(&cli.Command{
		Name:        "fake-data-fill",
		Description: "Create fake table and fill fake data",
		Action: func(ctx *cli.Context) error {
			return PgCreateFill(ctx.Context)
		},
	})
}

func printFValsForFakeData() *cli.Command {
	return pg.Wrap(&cli.Command{
		Name:        "fake-data-get-fvals",
		Description: "Prints filter values retrieved for an example fspec",
		Action: func(ctx *cli.Context) error {
			fval, err := pg.GetFilterValues(
				ctx.Context, "fake_user", UserFilterSpec, nil)
			if err != nil {
				return err
			}

			iox.PrintJSON(fval)
			return nil
		},
	})
}
