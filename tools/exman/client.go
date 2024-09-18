package main

import (
	"context"
	"fmt"

	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/proc"
)

type Client struct {
	client *httpx.Client
}

func NewClient(port uint32) *Client {
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	client := httpx.NewClient(url, "")
	return &Client{client: client}
}

func (c *Client) Exec(gtx context.Context, cmd *proc.CmdDesc) error {
	res := c.client.Build().Path("/api/v1/cmd").Post(gtx, cmd)
	if err := res.Close(); err != nil {
		return errx.Errf(err,
			"failed to add cmd '%s' to exec server", cmd.Name)
	}
	return nil
}

func (c *Client) List(gtx context.Context) ([]*proc.CmdInfo, error) {
	res := c.client.Build().Path("/api/v1/cmd").Get(gtx)
	cmds := make([]*proc.CmdInfo, 0, 20)
	if err := res.LoadClose(&cmds); err != nil {
		return nil, errx.Errf(err,
			"failed to get list of commands from exec server")
	}
	return cmds, nil
}

func (c *Client) Terminate(
	gtx context.Context, name string, force bool) error {
	res := c.client.Build().
		Path("/api/v1/cmds", name).
		QBool("force", force).
		Delete(gtx)
	if err := res.Close(); err != nil {
		return errx.Errf(err, "failed to terminate cmd '%s'", name)
	}
	return nil
}
