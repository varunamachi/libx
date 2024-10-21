package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/errx"
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

			// go func() {
			// 	<-gtx.Done()
			// 	fmt.Println("\nstopping server gracefully")
			// 	if err := server.server.Close(); err != nil {
			// 		log.Error().Err(err).Msg("failed to stop exec-server")
			// 	}

			// }()

			err := server.Start(gtx, "127.0.0.1", uint32(ctx.Uint("port")))

			if err != nil && !errors.Is(err, http.ErrServerClosed) {
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
				Required: false,
				Value:    "",
			},
			// &cli.StringFlag{
			// 	Name:     "cmd",
			// 	Usage:    "executable path or name without arguments",
			// 	Required: true,
			// },
			// &cli.StringFlag{
			// 	Name:     "args",
			// 	Usage:    "command arguments, space separated",
			// 	Required: true,
			// },
			&cli.BoolFlag{
				Name:  "fwd-env",
				Usage: "forward current env variables to server for this cmd",
				Value: true,
			},
			&cli.StringSliceFlag{
				Name:  "env",
				Usage: "Env vars in the form of key1=value1",
			},
			&cli.StringFlag{
				Name:  "cwd",
				Usage: "executables current working directory",
			},
		),
		Action: func(ctx *cli.Context) error {

			cwd := ctx.String("cwd")
			if cwd == "" {
				c, err := os.Getwd()
				if err != nil {
					return errx.Errf(err, "failed to commands cwd")
				}
				cwd = c
			}

			cmdName := ctx.Args().First()
			args := ctx.Args().Tail()

			envs := map[string]string{}
			for _, kv := range ctx.StringSlice("env") {
				comps := strings.Split(kv, ",")
				if len(comps) != 2 {
					log.Warn().Str("envVar", kv).
						Msg("invlid env variable given ignoring")
					continue
				}
				envs[comps[0]] = comps[1]
			}

			cmd := &proc.CmdDesc{
				Name:          ctx.String("name"),
				Path:          cmdName,
				Args:          args,
				Env:           envs,
				Cwd:           cwd,
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
				Required: false,
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force kill",
				Value: false,
			},
		),
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			if name == "" {
				name = ctx.Args().First()
			}
			err := client(ctx).Terminate(ctx.Context, name, ctx.Bool("force"))
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func stopAllCmd() *cli.Command {
	return &cli.Command{
		Name:        "stop-all",
		Usage:       "Stop a command running in exec-server",
		Description: "Stop a command running in exec-server",
		Flags: withServerFlags(
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force kill",
				Value: false,
			},
		),
		Action: func(ctx *cli.Context) error {
			err := client(ctx).TerminateAll(ctx.Context, ctx.Bool("force"))
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
