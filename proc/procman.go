package proc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

var (
	ErrProcessNotFound = errors.New("process not found")
	ErrCommandNotFound = errors.New("process not found")
)

type Manager struct {
	mutex sync.Mutex
	gtx   context.Context
	cmds  map[string]CmdEntry
}

func NewManager(gtx context.Context) *Manager {
	return &Manager{
		gtx:  gtx,
		cmds: map[string]CmdEntry{},
	}
}

func (man *Manager) Add(cdesc *CmdDesc) (int, error) {

	cmd := man.mkcmd(cdesc)

	if err := cmd.Start(); err != nil {
		return -1,
			errx.Errf(err,
				"failed to start command: %s - %s", cdesc.Name, cdesc.Path)
	}
	if cdesc.Name == "" {
		if cmd.Process != nil {
			cdesc.Name = fmt.Sprintf("%s-%d",
				filepath.Base(cmd.Path),
				cmd.Process.Pid)
		} else {
			cdesc.Name = filepath.Base(cmd.Path)
		}
		setName(cmd, cdesc.Name)
	}

	man.addToMap(cmd, cdesc)
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Fprintln(cmd.Stderr, err)
		}
		log.Info().Str("name", cdesc.Name).Msg("process exited")
		man.removeFromMap(cdesc.Name)
	}()

	return cmd.Process.Pid, nil
}

func (man *Manager) Terminate(name string, forceKill bool) error {
	cmd := man.Get(name)
	if cmd == nil {
		return errx.Errf(ErrCommandNotFound,
			"command with name '%s' does not exit")
	}

	if cmd.Process == nil {
		return errx.Errf(ErrProcessNotFound,
			"command '%s' does not have a associated process", name)
	}

	if err := cmd.Process.Kill(); err != nil {
		return errx.Errf(err, "failed to kill process '%d' for '%s'",
			cmd.Process.Pid, name)
	}
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

func (man *Manager) GetDesc(name string) *CmdDesc {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	excmd, found := man.cmds[name]
	if !found {
		return nil
	}
	return excmd.desc
}

func (man *Manager) List() []*CmdInfo {
	man.mutex.Lock()
	defer man.mutex.Unlock()

	out := make([]*CmdInfo, 0, len(man.cmds))
	for _, val := range man.cmds {
		out = append(out, &CmdInfo{
			Desc:    val.desc,
			Started: val.started,
		})
	}

	slices.SortFunc(out, func(a, b *CmdInfo) int {
		if a.Started.After(b.Started) {
			return 1
		}
		return -1
	})
	return out
}

func (man *Manager) addToMap(cmd *exec.Cmd, desc *CmdDesc) {
	man.mutex.Lock()
	defer man.mutex.Unlock()
	man.cmds[desc.Name] = CmdEntry{
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

	if desc.EnvsForwarded {
		// When envs are forwarded, we dont use server's envs
		cmd.Env = make([]string, 0, len(desc.Env))
	}

	for k, v := range desc.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	if desc.Cwd != "" {
		cmd.Dir = desc.Cwd
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = NewWriter(
		desc.Name,
		data.Qop[io.Writer](cmd.Stdout != nil, cmd.Stdout, os.Stdout),
		procNameStyle(),
	)
	cmd.Stderr = NewWriter(
		desc.Name,
		data.Qop[io.Writer](cmd.Stderr != nil, cmd.Stderr, os.Stderr),
		procNameStyle(),
	)
	return cmd
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

func setName(cmd *exec.Cmd, name string) {
	w, ok := cmd.Stdout.(*writer)
	if ok {
		w.SetName(name)
	}

	w, ok = cmd.Stderr.(*writer)
	if ok {
		w.SetName(name)
	}
}
