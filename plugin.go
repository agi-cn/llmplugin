package llmplugin

import "context"

type Plugin interface {
	Do(ctx context.Context, query string) (answer string, err error)

	GetName() string
	GetInputExample() string
	GetDesc() string
}

var _ Plugin = (*SimplePlugin)(nil)

type SimplePlugin struct {
	// Name of plugin.
	Name string

	// InputExample is the example of input.
	InputExample string

	// Desc is the description of plugin for LLM to understand what is the plugin and what for.
	Desc string

	// DoFunc is the handle function of plugin.
	DoFunc func(ctx context.Context, query string) (answer string, err error)
}

func (p SimplePlugin) GetName() string {
	return p.Name
}

func (p SimplePlugin) GetInputExample() string {
	return p.InputExample
}

func (p SimplePlugin) GetDesc() string {
	return p.Desc
}

func (p SimplePlugin) Do(ctx context.Context, query string) (answer string, err error) {
	return p.DoFunc(ctx, query)
}
