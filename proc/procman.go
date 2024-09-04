package proc

import (
	"context"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/varunamachi/libx/data"
)

type cmd struct {
	command *exec.Cmd
	desc    *CmdDesc
}

type Manager struct {
	gtx  context.Context
	cmds map[string]cmd
}

func (man *Manager) Add(cb *CmdDesc) (uint, error) {

	_ = man.mkcmd(cb)

	return 0, nil
}

func (man *Manager) mkcmd(desc *CmdDesc) *exec.Cmd {
	cmd := exec.CommandContext(man.gtx, desc.Path, desc.Args...)
	for k, v := range desc.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	if desc.Cwd != "" {
		cmd.Dir = desc.Cwd
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = NewWriter(
		desc.Name,
		data.Qop[io.Writer](cmd.Stdout == nil, cmd.Stdout, os.Stdout),
		procNameStyle(),
	)
	cmd.Stderr = NewWriter(
		desc.Name,
		data.Qop[io.Writer](cmd.Stderr == nil, cmd.Stderr, os.Stderr),
		procNameStyle(),
	)
	return nil

}

func procNameStyle() lipgloss.Style {

	// TODO - make this better
	// Choose light and dark colors and return them - probably not randomly
	color := lipgloss.Color(strconv.Itoa(rand.Intn(256)))
	return lipgloss.
		NewStyle().
		Foreground(color).
		Bold(true).
		Align(lipgloss.Left)
}
