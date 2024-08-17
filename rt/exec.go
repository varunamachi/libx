package rt

import (
	"io"
	"os"
	"os/exec"

	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
)

// TODO - add way to start command in background
// TODO - simple way to get output as string

type CmdBuilder struct {
	Cmd  string
	Env  map[string]string
	Args []string

	stdout io.Writer
	stdin  io.Reader
	stderr io.Writer
}

func NewCmdBuilder(cmd string) *CmdBuilder {
	return &CmdBuilder{
		Cmd: cmd,
	}
}

func (cb *CmdBuilder) WithEnv(name, value string) *CmdBuilder {
	if cb.Env == nil {
		cb.Env = map[string]string{}
	}
	cb.Env[name] = value
	return cb
}

func (cb *CmdBuilder) WithArgs(args ...string) *CmdBuilder {
	if cb.Args == nil {
		cb.Args = args
		return cb
	}
	cb.Args = append(cb.Args, args...)
	return cb
}

func (cb *CmdBuilder) WithOutput(out, err io.Writer) *CmdBuilder {
	cb.stdout, cb.stderr = out, err
	return cb
}

func (cb *CmdBuilder) WithCombinedOutput(out io.Writer) *CmdBuilder {
	cb.stdout, cb.stderr = out, out
	return cb
}

func (cb *CmdBuilder) WithInput(in io.Reader) *CmdBuilder {
	cb.stdin = in
	return cb
}

func (cb *CmdBuilder) Command() *exec.Cmd {
	cmd := exec.Command(cb.Cmd, cb.Args...)
	for k, v := range cb.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stdout = data.Qop(cb.stdout != nil, cb.stdout, io.Writer(os.Stdout))
	cmd.Stderr = data.Qop(cb.stderr != nil, cb.stderr, io.Writer(os.Stderr))
	cmd.Stdin = data.Qop(cb.stdin != nil, cb.stdin, io.Reader(os.Stdin))
	return cmd
}

func (cb *CmdBuilder) Run() error {
	cmd := cb.Command()
	if err := cmd.Run(); err != nil {
		return errx.Wrap(err)
	}
	return nil
}

func (cb *CmdBuilder) Start() (*os.Process, error) {
	cmd := cb.Command()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd.Process, nil
}
