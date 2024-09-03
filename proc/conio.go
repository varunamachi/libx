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
		_, err := fmt.Fprintf(cw.inner, "%s | %2s\n",
			cw.style.Render(cw.name), ln)
		if err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

func NewWriter(
	name string, target io.Writer, style lipgloss.Style) io.Writer {
	if len(name) < 10 {
		name = fmt.Sprintf("%10s", name)
	} else {
		name = fmt.Sprintf("%8s..", name[:8])
	}
	return &writer{
		name:  name,
		inner: target,
		style: style,
	}
}
