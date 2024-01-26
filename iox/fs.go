package iox

type WalkOptions struct {
}

func Walk(path *WalkOptions) error {
	return nil
}

func Glob(path string, filter func(string) bool) ([]string, error) {
	return nil, nil
}
