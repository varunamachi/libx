package proc

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type writer struct {
	name  string
	inner io.Writer
	style lipgloss.Style
}

func (cw *writer) Write(data []byte) (int, error) {
	strData := string(data)
	lines := strings.Split(strData, "\n")
	for _, ln := range lines {
		if ln == "" || strings.Contains(ln, "[sudo] password for") {
			continue
		}
		_, err := fmt.Fprintf(cw.inner, "%16s | %2s\n",
			cw.style.Render(cw.name), ln)
		if err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

func (cw *writer) SetName(name string) {
	if len(name) < 12 {
		cw.name = fmt.Sprintf("%12s", name)
	} else {
		cw.name = fmt.Sprintf("%110s..", name[:10])
	}
}

func NewWriter(
	name string, target io.Writer, style lipgloss.Style) io.Writer {
	if len(name) < 12 {
		name = fmt.Sprintf("%12s", name)
	} else {
		name = fmt.Sprintf("%10s..", name[:10])
	}
	return &writer{
		name:  name,
		inner: target,
		style: style,
	}
}
