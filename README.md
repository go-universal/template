# Template Library

![GitHub Tag](https://img.shields.io/github/v/tag/go-universal/template?sort=semver&label=version)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-universal/template.svg)](https://pkg.go.dev/github.com/go-universal/template)
[![License](https://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/go-universal/template/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-universal/template)](https://goreportcard.com/report/github.com/go-universal/template)
![Contributors](https://img.shields.io/github/contributors/go-universal/template)
![Issues](https://img.shields.io/github/issues/go-universal/template)

The `template` library is a Go package for rendering HTML templates with advanced features like partials, layouts, and custom template functions (pipes). It is designed to work seamlessly with a flexible filesystem.

## Installation

To install the library, use:

```bash
go get github.com/go-universal/template
```

## Features

- Efficient in-memory caching for templates.
- Flexible layout-based rendering.
- Support for global partial views.
- Robust error-safe rendering mechanism.

## Template Syntax

For layout template you can use `{{ view }}` function to render child view. All global partials template can accessed by `@partials/path/to/file` or `template-name`.

**NOTE**: Use `include` function to import instead of builtin `template` function to prevent errors.

**NOTE**: In `Render`, you can load multiple partial views. The first arguments must be the layout, and if the view has a partial and no layout, an empty string must be passed as the first argument.

### Builtin functions

- `{{ view }}`: render child template in layout. If used in non-layout template generate error!
- `{{ exists "template name or path" }}`: check if template name or path exists.
- `{{ include "template name or path" (optional data) }}`: includes and executes a template with the given name or path and data if exists.
- `{{ require "template name or path" (optional data) }}`: includes and executes a template with the given name or path and data or returning an error if the template does not exist.

## Usage

### Basic Example

```html
<!-- parts/header.tpl -->
{{ define "site-header" }}
<header>...</header>
{{ end }}

<!-- parts/sub/footer.tpl -->
<footer>...</footer>

<!-- pages/home.tpl -->
<section>
    <h1>Home Page</h1>
    <p>{{ .Title }}</p>
</section>
{{ define "title" }}Home Page{{ end }}

<!-- layout.tpl -->
<html>
    <head>
        {{ if exists "title" }}
            <title>{{ include "title" }}</title>
        {{ else }}
            <title>My App</title>
        {{ end }}
    </head>
    <body>
        {{- require "site-header" . }}
        {{- view }}
        {{- include "@partials/sub/footer" }}
    </body>
</html>
```

```go
package main

import (
    "os"
    "github.com/go-universal/fs"
    "github.com/go-universal/template"
)

func main() {
    source := fs.NewDir("./views")
    tpl := template.New(source, template.WithPartials("parts"))

    err := tpl.Load()
    if err != nil {
        panic(err)
    }


    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      data := template.Ctx().Add("Title", "Hello, World!")
      err = tpl.Render(w, "pages/home", data, "layout")
      if err != nil {
        // Render error page with no layout and two partial views
        tpl.Render(w, "pages/errors", nil, "", "pages/err-partial/500", "pages/err-partial/contact")
      }
    })

}
```

### Custom Options

```go
tpl := template.New(fs,
    template.WithRoot("."),
    template.WithPartials("partials"),
    template.WithExtension(".tpl"),
    template.WithDelimeters("{{", "}}"),
    template.WithEnv(true),
    template.WithCache(),
    template.WithUUIDPipe(),
    template.WithTernaryPipe(),
)
```

## API

### Template Interface

```go
type Template interface {
    Load() error
    Render(w io.Writer, view string, data interface{}, layouts ...string) error
    Compile(name, layout string, data any) ([]byte, error)
}
```

### Options

- `WithRoot(root string) Options`: Sets the root directory for templates.
- `WithPartials(path string) Options`: Sets the directory for partial templates.
- `WithExtension(ext string) Options`: Sets the file extension for templates.
- `WithDelimeters(left, right string) Options`: Sets the delimiters for template tags.
- `WithEnv(isDev bool) Options`: Sets the environment mode (development or production).
- `WithCache() Options`: Enables template caching.
- `WithPipes(name string, fn any) Options`: Registers a custom function (pipe) for templates.

### Context

Helper struct to pass data to template.

```go
type Context struct {
    data map[string]any
}

func Ctx() *Context
func ToCtx(v any) *Context
func (ctx *Context) Add(k string, v any) *Context
func (ctx *Context) Map() map[string]any
```

### Custom Pipes

- `WithUUIDPipe() Options`: Adds a UUID generation pipe.
- `WithTernaryPipe() Options`: Adds a ternary operation pipe.
- `WithNumberFmtPipe() Options`: Adds a number formatting pipe.
- `WithRegexpFmtPipe() Options`: Adds a regular expression formatting pipe.
- `WithJSONPipe() Options`: Adds a JSON formatting pipe.
- `WithDictPipe() Options`: Adds a dictionary creation pipe.
- `WithIsSetPipe() Options`: Adds a pipe to check if a value is set.
- `WithAlterPipe() Options`: Adds a pipe to alter a value.
- `WithDeepAlterPipe() Options`: Adds a pipe to deeply alter a value.
- `WithBrPipe() Options`: Adds a pipe to convert `\n` to `<br>`.

## License

This library is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.
