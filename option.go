package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strings"

	"github.com/go-universal/utils"
	"github.com/google/uuid"
)

type option struct {
	root       string
	partials   string
	extension  string
	leftDelim  string
	rightDelim string
	Dev        bool
	Cache      bool
	Pipes      template.FuncMap
}

// Options represents a configuration option for the Template.
type Options func(*option)

// WithRoot sets the root directory for templates. Default is ".".
func WithRoot(root string) Options {
	root = normalizePath(root)
	return func(opt *option) {
		opt.root = root + "/"
		if root == "" {
			opt.root = "."
		}
	}
}

// WithPartials sets the partials path for templates.
func WithPartials(path string) Options {
	path = normalizePath(path)
	return func(opt *option) {
		if path != "" && path != "." {
			opt.partials = path + "/"
		}
	}
}

// WithExtension sets the file extension for templates. Default is ".tpl".
func WithExtension(ext string) Options {
	ext = strings.TrimSpace(ext)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return func(opt *option) {
		if ext != "" {
			opt.extension = ext
		}
	}
}

// WithDelimeters sets custom delimiters for templates. Default is "{{" and "}}".
func WithDelimeters(left, right string) Options {
	left, right = strings.TrimSpace(left), strings.TrimSpace(right)
	return func(opt *option) {
		if left != "" && right != "" {
			opt.leftDelim, opt.rightDelim = left, right
		}
	}
}

// WithEnv sets the environment to development or production mode.
func WithEnv(isDev bool) Options {
	return func(opt *option) {
		opt.Dev = isDev
	}
}

// WithCache enables caching for templates. Disabled by default.
func WithCache() Options {
	return func(opt *option) {
		opt.Cache = true
	}
}

// WithPipes registers a custom function for use in templates.
func WithPipes(name string, fn any) Options {
	name = strings.TrimSpace(name)
	return func(opt *option) {
		if name != "" && fn != nil {
			opt.Pipes[name] = fn
		}
	}
}

// WithUUIDPipe adds a "uuid" pipe to generate UUID strings.
func WithUUIDPipe() Options {
	return func(opt *option) {
		opt.Pipes["uuid"] = func() string {
			return uuid.NewString()
		}
	}
}

// WithTernaryPipe adds an "iif" pipe for ternary operations.
func WithTernaryPipe() Options {
	return func(opt *option) {
		opt.Pipes["iif"] = func(cond bool, y, n any) any {
			if cond {
				return y
			}
			return n
		}
	}
}

// WithNumberFmtPipe adds a "numberFmt" pipe to format numbers.
func WithNumberFmtPipe() Options {
	return func(opt *option) {
		opt.Pipes["numberFmt"] = func(layout string, v ...any) string {
			return utils.FormatNumber(layout, v...)
		}
	}
}

// WithRegexpFmtPipe adds a "regexpFmt" pipe to format strings using regex.
func WithRegexpFmtPipe() Options {
	return func(opt *option) {
		opt.Pipes["regexpFmt"] = func(data, pattern, repl string) (string, error) {
			return utils.FormatRx(data, pattern, repl)
		}
	}
}

// WithJSONPipe adds a "toJson" pipe to convert data to JSON strings.
func WithJSONPipe() Options {
	return func(opt *option) {
		opt.Pipes["toJson"] = func(data any) (string, error) {
			res, err := json.Marshal(data)
			if err != nil {
				return "", err
			}
			return string(res), nil
		}
	}
}

// WithDictPipe adds a "dict" pipe to create a map from key-value pairs.
//
// code block:
//
//	{{ $userGlobal := toJson
//				"name" .User.Name
//				"family" .User.Family
//				"email" .User.Email
//	}}
func WithDictPipe() Options {
	return func(opt *option) {
		opt.Pipes["dict"] = func(kv ...any) (map[string]any, error) {
			if len(kv)%2 != 0 {
				return nil, fmt.Errorf("invalid number of arguments for dict")
			}
			dict := make(map[string]any)
			for i := 0; i < len(kv); i += 2 {
				key, ok := kv[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = kv[i+1]
			}
			return dict, nil
		}
	}
}

// WithIsSetPipe adds an "isSet" pipe to check if a field exists in a map.
//
// code block:
//
//	{{ if (isSet dataMap "title") }}
//		<h1>{{ dataMap.title }}</h1>
//	{{ else }}
//		<h1>Unknown</h1>
//	{{ end }}
func WithIsSetPipe() Options {
	return func(opt *option) {
		opt.Pipes["isSet"] = func(data map[string]any, field string) bool {
			_, ok := data[field]
			return ok
		}
	}
}

// WithAlterPipe adds an "alter" pipe to return an alternative value if the original is nil.
//
// code block:
//
//	{{ $safeTitle := alter .data.meta.Title "Greeting" }}
func WithAlterPipe() Options {
	return func(opt *option) {
		opt.Pipes["alter"] = func(val, alt any) any {
			if val == nil {
				return alt
			}
			return val
		}
	}
}

// WithDeepAlterPipe adds a "deepAlter" pipe to handle nil or zero values.
//
// code block:
//
//	{{ $safeTitle := deepAlter .data.meta.Title "Greeting" }}
func WithDeepAlterPipe() Options {
	return func(opt *option) {
		opt.Pipes["deepAlter"] = func(val, alt any) any {
			if val == nil {
				return alt
			}
			v := reflect.ValueOf(val)
			switch v.Kind() {
			case reflect.String, reflect.Slice, reflect.Map, reflect.Chan:
				if v.Len() == 0 {
					return alt
				}
			case reflect.Ptr, reflect.Interface:
				if v.IsNil() {
					return alt
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				if v.IsZero() {
					return alt
				}
			}
			return val
		}
	}
}

// WithBrPipe adds a "br" pipe to replace newlines with HTML line breaks.
//
// code block:
//
//	{{ $out := br .comment }}
func WithBrPipe() Options {
	return func(opt *option) {
		opt.Pipes["br"] = func(text string) template.HTML {
			escaped := template.HTMLEscapeString(text)
			return template.HTML(strings.ReplaceAll(escaped, "\n", "<br/>"))
		}
	}
}
