package template

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

// viewPipe registers a custom "view" function for rendering a child template
// inside a layout template. It returns an error if the child template fails
// to render or if "view" is called from a non-layout template.
func viewPipe(t *template.Template, data []byte) {
	t.Funcs(map[string]any{
		"view": func() (template.HTML, error) {
			if data == nil {
				return "", errors.New("layout template called without view")
			}
			return template.HTML(data), nil
		},
	})
}

// existsPipe registers a custom "exists" function to the template engine.
// The "exists" function checks if a template with the given name exists.
func existsPipe(t *template.Template) {
	t.Funcs(map[string]any{
		"exists": func(name string) bool {
			return t.Lookup(name) != nil
		},
	})
}

// includePipe registers a custom "include" function to the template engine.
// The "include" function includes and executes a template with the given name.
// If the template does not exist, it returns an empty string without error.
func includePipe(t *template.Template) {
	t.Funcs(map[string]any{
		"include": func(name string, data ...any) (template.HTML, error) {
			tpl := t.Lookup(name)
			if tpl == nil {
				return "", nil
			}

			var v any
			if len(data) > 0 {
				v = data[0]
			}

			var buf bytes.Buffer
			if err := tpl.Execute(&buf, underlyingValue(v)); err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
	})
}

// requirePipe registers a custom "require" function to the template engine.
// The "require" function includes and executes a template with the given name.
// If the template does not exist, it returns an error.
func requirePipe(t *template.Template) {
	t.Funcs(map[string]any{
		"require": func(name string, data ...any) (template.HTML, error) {
			tpl := t.Lookup(name)
			if tpl == nil {
				return "", fmt.Errorf("template %s does not exist", name)
			}

			var v any
			if len(data) > 0 {
				v = data[0]
			}

			var buf bytes.Buffer
			if err := tpl.Execute(&buf, underlyingValue(v)); err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
	})
}
