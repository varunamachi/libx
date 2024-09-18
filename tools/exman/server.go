package main

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/proc"
)

type Server struct {
	server *httpx.Server
	man    *proc.Manager
}

func (s *Server) Start(gtx context.Context, bindIp string, port uint32) error {
	s.server = httpx.NewServer(os.Stdout, nil)

	s.server.WithAPIs(
		s.executeEp(),
		s.terminateEp(),
		s.listEp(),
	)

	if err := s.server.StartContext(gtx, port); err != nil {
		if err != http.ErrServerClosed {
			return errx.Wrap(err)
		}
	}
	return nil
}

func (s *Server) executeEp() *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var desc proc.CmdDesc
		if err := etx.Bind(&desc); err != nil {
			return errx.BadReqX(err, "failed to read command from request")
		}

		if _, err := s.man.Add(&desc); err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, data.M{
			"started": desc.Name,
		})
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

func (s *Server) terminateEp() *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		pgr := httpx.NewParamGetter(etx)
		name := pgr.Str("name")
		force := pgr.QueryBoolOr("force", false)
		if err := s.man.Terminate(name, force); err != nil {
			return err
		}

		return httpx.SendJSON(etx, data.M{
			"deleted": name,
		})
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

func (s *Server) listEp() *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		list := s.man.List()
		return httpx.SendJSON(etx, list)
	}

	return &httpx.Endpoint{
		Method:   echo.GET,
		Path:     "cmd",
		Category: "cmd-exec",
		Desc:     "Get a list of managed commands",
		Version:  "v1",
		Handler:  handler,
	}
}
