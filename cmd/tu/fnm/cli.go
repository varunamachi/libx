package fnm

import (
	"errors"

	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/iox"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		pg.Wrap(&cli.Command{
			Name:        "fake-data",
			Description: "Commands related managing fake data for testing",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "data-type",
					Usage: "Name of the dataType/table to act on. If empty " +
						"all or any will be taken based on context",
					Value: "",
				},
			},
			Subcommands: []*cli.Command{
				fillData(),
				printFilterValues(),
				createAndApplyRandomFilter(),
			},
		}),
	}
}

func fillData() *cli.Command {
	return &cli.Command{
		Name: "fill",
		Description: "Create fake table and fill fake data, " +
			"ignores data type flag",
		Action: func(ctx *cli.Context) error {
			return PgCreateFill(ctx.Context)
		},
	}
}

func printFilterValues() *cli.Command {
	return &cli.Command{
		Name:        "get-fvals",
		Description: "Prints filter values retrieved for an example fspec",
		Action: func(ctx *cli.Context) error {
			fval, err := pg.GetFilterValues(
				ctx.Context, ctx.String("data-type"), UserFilterSpec, nil)
			if err != nil {
				return err
			}

			iox.PrintJSON(fval)
			return nil
		},
	}
}

func createAndApplyRandomFilter() *cli.Command {
	return &cli.Command{
		Name:        "get-random-filtered",
		Description: "Creates a random filter and applies it on given table",
		Flags:       []cli.Flag{},
		Action: func(ctx *cli.Context) error {

			var out any

			switch ctx.String("data-type") {
			case "fake_user":
				a := make([]*FkUser, 0, 100)
				out = &a
			case "fake_item":
				b := make([]*FkItem, 0, 100)
				out = &b
			default:
				return errx.Errf(errors.New("invalid data type"),
					"invalid data type '%s'", ctx.String("data-type"))
			}

			err := GetDataForRandomFilter(
				ctx.Context,
				ctx.String("data-type"),
				out)
			if err != nil {
				return err
			}

			iox.PrintJSON(out)
			return nil
		},
	}
}
