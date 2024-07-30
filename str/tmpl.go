package str

import (
	"bytes"
	ht "html/template"
	"io"
	tt "text/template"

	"github.com/varunamachi/libx/errx"
)

type TemplateDesc struct {
	Template string
	Data     map[string]interface{}
	Funcs    map[string]any
	Html     bool
}

func SimpleTemplateExpand(td *TemplateDesc) (string, error) {
	buf := bytes.Buffer{}
	expander := TextTemplateExpand
	if td.Html {
		expander = HtmlTemplateExpand
	}
	if err := expander(td, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TextTemplateExpand(td *TemplateDesc, writer io.Writer) error {

	t, err := tt.New("tmpl").Parse(td.Template)
	if err != nil {
		return errx.Errf(err, "failed to parse html template")
	}

	if len(td.Funcs) != 0 {
		t.Funcs(td.Funcs)
	}

	if err := t.Execute(writer, td.Data); err != nil {
		return errx.Errf(err, "failed to execute html template")
	}
	return nil
}

func HtmlTemplateExpand(td *TemplateDesc, writer io.Writer) error {

	t, err := ht.New("tmpl").Parse(td.Template)
	if err != nil {
		return errx.Errf(err, "failed to parse html template")
	}

	if len(td.Funcs) != 0 {
		t.Funcs(td.Funcs)
	}

	if err := t.Execute(writer, td.Data); err != nil {
		return errx.Errf(err, "failed to execute html template")
	}

	return nil
}
