package httpx

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/iox"
)

func WithClientFlags(
	withAuth bool,
	envPrefix string,
	flags ...cli.Flag) []cli.Flag {

	prefix := data.Qop(envPrefix == "", "LIBX", envPrefix)
	ev := func(envVar string) []string {
		return []string{prefix + "_" + envVar}
	}

	flags = append(flags,
		&cli.StringFlag{
			Name:     "url",
			Usage:    "URL of the server",
			EnvVars:  ev("REMOTE_URL"),
			Required: true,
		},
		&cli.BoolFlag{
			Name: "ignore-cert-errors",
			Usage: "Ignore certificate errors while connecting to a HTTPS " +
				"service",
			Value:   false,
			EnvVars: ev("CLIENT_IGNORE_CERT_ERR"),
		},
		&cli.IntFlag{
			Name:    "timeout-secs",
			Usage:   "Time out in seconds",
			Value:   20,
			EnvVars: ev("CLIENT_TIMEOUT_SECS"),
		},
	)
	if withAuth {
		flags = append(flags,
			&cli.StringFlag{
				Name:     "user-id",
				Usage:    "User present in the remote service",
				Required: false,
				EnvVars:  ev("CLIENT_USER_ID"),
			},
			&cli.StringFlag{
				Name: "password",
				Usage: "Password for the user, only use for development " +
					"purposes",
				Required: false,
				Hidden:   true,
				EnvVars:  ev("CLIENT_PASSWORD"),
			},
		)
	}

	return flags
}

func CreateClient(ctx *cli.Context) (*Client, AuthData, error) {
	host := ctx.String("host")
	ignCertErrs := ctx.Bool("ignore-cert-errors")
	timeOut := ctx.Int("timeout-secs")

	tp := DefaultTransport()
	if ignCertErrs {
		tp.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: ignCertErrs,
		}
	}
	client := NewCustomClient(host, "", tp, time.Duration(timeOut)*time.Second)

	userId := ctx.String("user-id")
	if userId == "" {
		return client, nil, nil
	}

	password := ctx.String("password")
	if password == "" {
		password = iox.AskPassword(fmt.Sprintf("Password for '%s'", userId))
	}

	authData := AuthData{
		"userId":   userId,
		"password": password,
	}

	return client, authData, nil
}
