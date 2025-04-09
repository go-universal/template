package template

// Context represents a collection of key-value pairs for template data.
type Context struct {
	data map[string]any
}

// Ctx creates and returns a new empty Context instance.
func Ctx() *Context {
	return &Context{
		data: make(map[string]any),
	}
}

// ToContext converts the given value to a Context instance.
// If the value is not a valid map or Context, it returns a new empty Context.
func ToContext(v any) *Context {
	switch val := v.(type) {
	case map[string]any:
		return &Context{data: val}
	case Context:
		return &val
	case *Context:
		return val
	default:
		return Ctx()
	}
}

// Add inserts a key-value pair into the Context.
// If the key is empty, the operation is ignored.
func (ctx *Context) Add(key string, value any) *Context {
	if key != "" {
		ctx.data[key] = value
	}
	return ctx
}

// Data returns the underlying map of the Context.
func (ctx *Context) Data() map[string]any {
	return ctx.data
}
