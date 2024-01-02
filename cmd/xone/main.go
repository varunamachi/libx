package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/rt"
	"github.com/varunamachi/libx/testutils/fake"
)

// var test2 = pg.NewSelectorGenerator().SelectorX(&data.CommonParams{
// 	Filter: &data.Filter{
// 		// Bools: map[string]any{
// 		// 	"boolOne":   true,
// 		// 	"boolTwo":   false,
// 		// 	"boolThree": nil,
// 		// },
// 		Props: map[string]*data.Matcher{
// 			"variety_en": {
// 				Fields: []any{"abc", "def", "ghi"},
// 			},
// 			"unit_en": {
// 				Fields: []any{"jkl", "mno", "pqr"},
// 				Invert: true,
// 			},
// 		},
// 		Lists: map[string]*data.Matcher{
// 			"listOne": {
// 				Fields: []any{"Labc", "Ldef", "Lghi"},
// 			},
// 			"listTwo": {
// 				Fields: []any{"Ljkl", "Lmno", "Lpqr"},
// 				Invert: true,
// 			},
// 		},
// 		Searches: map[string]*data.Matcher{
// 			"searchOne": {
// 				Fields: []any{"Sabc", "Sdef", "Sghi"},
// 			},
// 			"searchTwo": {
// 				Fields: []any{"Sjkl", "Smno", "Spqr"},
// 				Invert: true,
// 			},
// 		},
// 		Constants: map[string]*data.Matcher{
// 			"constantOne": {
// 				Fields: []any{"Cabc", "Cdef", "Cghi"},
// 			},
// 			"constantTwo": {
// 				Fields: []any{"Cjkl", "Cmno", "Cpqr"},
// 				Invert: true,
// 			},
// 		},
// 		Dates: map[string]*data.DateRangeMatcher{
// 			"dateOne": {
// 				DateRange: data.DateRange{
// 					From: time.Now().AddDate(0, -1, 0),
// 					To:   time.Now(),
// 				},
// 				Invert: false,
// 			},
// 			"dateTwo": {
// 				DateRange: data.DateRange{
// 					From: time.Now().AddDate(0, 0, -7),
// 					To:   time.Now().AddDate(0, 0, -2),
// 				},
// 				Invert: true,
// 			},
// 		},
// 		Ranges: map[string]*data.RangeMatcher{
// 			"rangeOne": {
// 				NumberRange: data.NumberRange{
// 					From: 0,
// 					To:   100,
// 				},
// 				Invert: false,
// 			},
// 			"rangeTwo": {
// 				NumberRange: data.NumberRange{
// 					From: 50,
// 					To:   60,
// 				},
// 				Invert: true,
// 			},
// 		},
// 	},
// 	Page:           0,
// 	PageSize:       100,
// 	Sort:           "www",
// 	SortDescending: true,
// })

func main() {
	gtx, cancel := rt.Gtx()
	defer cancel()

	app := cli.NewApp()
	app.Commands = append(app.Commands, fake.FillCmd())

	if err := app.RunContext(gtx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func queryGen() {
	sel := pg.NewSelectorGenerator().SelectorX(&data.CommonParams{
		Filter: &data.Filter{
			Bools: map[string]any{
				"boolOne":   true,
				"boolTwo":   false,
				"boolThree": nil,
			},
			Props: map[string]*data.Matcher{
				"propOne": {
					Fields: []any{"abc", "def", "ghi"},
				},
				"propTwo": {
					Fields: []any{"jkl", "mno", "pqr"},
					Invert: true,
				},
			},
			Lists: map[string]*data.Matcher{
				"listOne": {
					Fields: []any{"Labc", "Ldef", "Lghi"},
				},
				"listTwo": {
					Fields: []any{"Ljkl", "Lmno", "Lpqr"},
					Invert: true,
				},
			},
			Searches: map[string]*data.Matcher{
				"searchOne": {
					Fields: []any{"Sabc", "Sdef", "Sghi"},
				},
				"searchTwo": {
					Fields: []any{"Sjkl", "Smno", "Spqr"},
					Invert: true,
				},
			},
			Constants: map[string]*data.Matcher{
				"constantOne": {
					Fields: []any{"Cabc", "Cdef", "Cghi"},
				},
				"constantTwo": {
					Fields: []any{"Cjkl", "Cmno", "Cpqr"},
					Invert: true,
				},
			},
			Dates: map[string]*data.DateRangeMatcher{
				"dateOne": {
					DateRange: data.DateRange{
						From: time.Now().AddDate(0, -1, 0),
						To:   time.Now(),
					},
					Invert: false,
				},
				"dateTwo": {
					DateRange: data.DateRange{
						From: time.Now().AddDate(0, 0, -7),
						To:   time.Now().AddDate(0, 0, -2),
					},
					Invert: true,
				},
			},
			Ranges: map[string]*data.RangeMatcher{
				"rangeOne": {
					NumberRange: data.NumberRange{
						From: 0,
						To:   100,
					},
					Invert: false,
				},
				"rangeTwo": {
					NumberRange: data.NumberRange{
						From: 50,
						To:   60,
					},
					Invert: true,
				},
			},
		},
		Page:           0,
		PageSize:       100,
		Sort:           "www",
		SortDescending: true,
	})

	sq := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("*").
		From("price_stat").
		Where(sel.QueryFragment, sel.Args...)
	query, args, err := sq.ToSql()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build sql query")
	}

	fmt.Println(query)
	fmt.Println()
	fmt.Println(args)

}
