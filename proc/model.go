package proc

import (
	"os/exec"
	"time"
)

type CmdDesc struct {
	Name          string
	Path          string
	Args          []string
	Env           map[string]string
	Cwd           string
	EnvsForwarded bool
}

type CmdEntry struct {
	command *exec.Cmd
	desc    *CmdDesc
	started time.Time
}

type CmdInfo struct {
	Desc    *CmdDesc
	Started time.Time
	PID     int
}
