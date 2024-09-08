package proc

import (
	"context"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

type execdCmd struct {
	command *exec.Cmd
	desc    *CmdDesc
}

type Manager struct {
	mutex sync.Mutex
	gtx   context.Context
	cmds  map[string]execdCmd
}

func (man *Manager) Add(cdesc *CmdDesc) (int, error) {

	cmd := man.mkcmd(cdesc)

	if err := cmd.Start(); err != nil {
		return -1,
			errx.Errf(err,
				"failed to start command: %s - %s", cdesc.Name, cdesc.Path)
	}
	man.addToMap(cmd, cdesc)
	go func() {
		cmd.Wait()
		man.removeFromMap(cdesc.Name)
	}()

	return cmd.Process.Pid, nil
}

func (man *Manager) IsRunning(name string) (bool, error) {
	return false, nil
}

func (man *Manager) Terminate(name string, forceKill bool) error {
	return nil
}

func (man *Manager) Get(name string) *exec.Cmd {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	excmd, found := man.cmds[name]
	if !found {
		return nil
	}
	return excmd.command
}

func (man *Manager) addToMap(cmd *exec.Cmd, desc *CmdDesc) {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	man.cmds[desc.Name] = execdCmd{
		command: cmd,
		desc:    desc,
	}
}

func (man *Manager) removeFromMap(name string) {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	delete(man.cmds, name)
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
