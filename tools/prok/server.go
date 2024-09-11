package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/proc"
)

type Server struct {
	server *httpx.Server
	man    *proc.Manager
}

func (s *Server) Start(bindIp string, port int) error {
	s.server = httpx.NewServer(os.Stdout, nil)

	if err := app.Serve(uint32(ctx.Uint("port"))); err != nil {
		if err != http.ErrServerClosed {
			return errx.Wrap(err)
		}
	}
	return nil
}

func executeEp() *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		return nil
	}

	return &httpx.Endpoint{
		Method:   echo.POST,
		Path:     "cmd",
		Category: "cmd-exec",
		Desc:     "Add a command to command/proc manager",
		Version:  "v1",
		Handler:  handler,
	}
}

func terminateEp() *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		return nil
	}

	return &httpx.Endpoint{
		Method:   echo.DELETE,
		Path:     "cmd/:name",
		Category: "cmd-exec",
		Desc:     "Kill a process that was started using command/proc manager",
		Version:  "v1",
		Handler:  handler,
	}
}
