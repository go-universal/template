package template

import (
	"path/filepath"
	"regexp"
	"strings"
)

// normalizePath joins and cleans paths, then converts them to use slashes as separators.
func normalizePath(paths ...string) string {
	return filepath.ToSlash(filepath.Clean(filepath.Join(paths...)))
}

// toName converts a file path to a name by removing the root and extension.
func toName(path, root, ext string) string {
	if path == "" {
		return ""
	}
	path = strings.TrimPrefix(path, root)
	path = strings.TrimSuffix(path, ext)
	return normalizePath(path)
}

// toPath converts a name to a file path by appending the root and extension.
func toPath(name, root, ext string) string {
	if name == "" {
		return ""
	}
	name = strings.TrimPrefix(name, root)
	name = strings.TrimSuffix(name, ext)
	return normalizePath(root, name+ext)
}

// toKey generates a unique key by concatenating multiple view names with a colon separator.
func toKey(views ...string) string {
	var res strings.Builder
	for _, v := range views {
		if v != "" {
			if res.Len() > 0 {
				res.WriteString(":")
			}
			res.WriteString(v)
		}
	}
	return res.String()
}

// underlyingValue extracts the underlying data from a Context type.
func underlyingValue(v any) any {
	switch val := v.(type) {
	case Context:
		return val.data
	case *Context:
		return val.data
	default:
		return v
	}
}

// extPattern creates a regular expression pattern to match paths with a specific extension.
func extPattern(path, ext string) string {
	if path == "" {
		return ".*" + regexp.QuoteMeta(ext)
	}
	return "^" + regexp.QuoteMeta(path) + ".*" + regexp.QuoteMeta(ext)
}
