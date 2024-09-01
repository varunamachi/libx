package proc

import "os/exec"

type ManagerConfig struct {
}

type Manager struct {
}

func (man *Manager) Add(name string, cb exec.Cmd) (uint, error) {
	return 0, nil
}
