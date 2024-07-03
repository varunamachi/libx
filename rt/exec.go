package rt

import "io"

type CmdBuilder struct {
	Env    map[string]any
	Args   []string
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer
}
