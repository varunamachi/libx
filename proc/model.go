package proc

type CmdDesc struct {
	Name string
	Path string
	Args []string
	Env  map[string]string
	Cwd  string
}
