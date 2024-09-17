package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/proc"
)

var GitTag = "--"
var GitHash = "--"
var GitBranch = "--"
var BuildTime = "--"
var BuildHost = "--"
var BuildUser = "--"

var bi = libx.BuildInfo{
	GitTag:    GitTag,
	GitHash:   GitHash,
	GitBranch: GitBranch,
	BuildTime: BuildTime,
	BuildHost: BuildHost,
	BuildUser: BuildUser,
}

func infoCmd() *cli.Command {
	return &cli.Command{
		Name:        "build-info",
		Usage:       "Prints build info for the tool",
		Description: "Prints build info for the tool",
		Flags:       withServerFlags(),
		Action: func(ctx *cli.Context) error {

			fmt.Println("GitTag:    ", bi.GitTag)
			fmt.Println("GitHash:   ", bi.GitHash)
			fmt.Println("GitBranch: ", bi.GitBranch)
			fmt.Println("BuildTime: ", bi.BuildTime)
			fmt.Println("BuildHost: ", bi.BuildHost)
			fmt.Println("BuildUser: ", bi.BuildUser)

			return nil
		},
	}
}

func serveCmd(gtx context.Context) *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Usage:       "Start the process manager server",
		Description: "Start the process manager server",
		Flags:       withServerFlags(),
		Action: func(ctx *cli.Context) error {
			server := Server{
				server: httpx.NewServer(os.Stdout, nil),
				man:    proc.NewManager(gtx),
			}

			err := server.Start("127.0.0.1", uint32(ctx.Uint("port")))
			if err != nil {
				return err
			}
			return nil
		},
	}
}

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
		Flags:       withServerFlags(),
		Action: func(ctx *cli.Context) error {
			list, err := client(ctx).List(ctx.Context)
			if err != nil {
				return err
			}

			// TODO - add advanced printing
			for _, ci := range list {
				fmt.Println(ci.Desc.Name)
			}

			return nil
		},
	}
}

func stopCmd() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       "Stop a command running in exec-server",
		Description: "Stop a command running in exec-server",
		Flags: withServerFlags(
			&cli.StringFlag{
				Name:     "name",
				Usage:    "name for this instance of the command",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force kill",
				Value: false,
			},
		),
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			err := client(ctx).Terminate(ctx.Context, name, ctx.Bool("force"))
			if err != nil {
				return err
			}
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
