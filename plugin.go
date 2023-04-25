package llmplugin

import "context"

type Plugin interface {
	Do(ctx context.Context, query string) (answer string, err error)

	GetName() string

	GetDesc() string
}

var _ Plugin = (*SimplePlugin)(nil)

type SimplePlugin struct {
	Name   string
	Desc   string
	DoFunc func(ctx context.Context, query string) (answer string, err error)
}

func (p SimplePlugin) GetName() string {
	return p.Name
}

func (p SimplePlugin) GetDesc() string {
	return p.Desc
}

func (p SimplePlugin) Do(ctx context.Context, query string) (answer string, err error) {
	return p.DoFunc(ctx, query)
}
