package main

import (
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/proc"
)

func execCmd() *cli.Command {
	return &cli.Command{
		Name:        "exec",
		Usage:       "Execute a command on exec-server",
		Description: "Execute a command on exec-server",
		Flags: withServerFlags(
			&cli.StringFlag{
				Name:     "name",
				Usage:    "name for this instance of the command",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "cmd",
				Usage:    "executable path or name without arguments",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "args",
				Usage:    "command arguments, space separated",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "fwd-env",
				Usage: "forward current env variables to server for this cmd",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "env",
				Usage: "Env vars in the form of key1=value1,key2=value2",
				Value: true,
			},
		),
		Action: func(ctx *cli.Context) error {
			cmd := &proc.CmdDesc{
				Name:          ctx.String("name"),
				Path:          ctx.String("cmd"),
				Args:          strings.Fields(ctx.String("args")),
				Env:           map[string]string{},
				Cwd:           "",
				EnvsForwarded: ctx.Bool("fwd-env"),
			}
			return client(ctx).Exec(ctx.Context, cmd)
		},
	}
}

func listCmd() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List commands running in exec-server",
		Description: "List commands running in exec-server",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}

func stopCmd() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       "Stop a command running in exec-server",
		Description: "Stop a command running in exec-server",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}

func withServerFlags(flags ...cli.Flag) []cli.Flag {
	flags = append(flags, &cli.UintFlag{
		Name:  "port",
		Value: 12012,
		Usage: "port number at which exec server is running",
	})
	return flags
}

func client(ctx *cli.Context) *Client {
	return NewClient(uint32(ctx.Uint("port")))
}
