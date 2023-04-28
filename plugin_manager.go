package llmplugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/agi-cn/llmplugin/llm"

	"github.com/sirupsen/logrus"
)

type PluginContext struct {
	Plugin

	// Input for handle function of plugin.
	Input string
}

type PluginManager struct {
	llmer llm.LLMer

	// plugins <key:name, value:Plugin>
	plugins map[string]Plugin
}

type PluginManagerOpt func(manager *PluginManager)

// WithPlugin enable one plugin.
func WithPlugin(p Plugin) PluginManagerOpt {

	return func(manager *PluginManager) {
		name := strings.ToLower(p.GetName())
		if _, ok := manager.plugins[name]; !ok {
			manager.plugins[name] = p
		}
	}
}

// WithPlugins enable multiple plugins.
func WithPlugins(plugins []Plugin) PluginManagerOpt {

	return func(manager *PluginManager) {

		for _, p := range plugins {
			opt := WithPlugin(p)
			opt(manager)
		}
	}
}

// NewPluginManager create plugin manager.
func NewPluginManager(llmer llm.LLMer, opts ...PluginManagerOpt) *PluginManager {

	manager := &PluginManager{
		llmer:   llmer,
		plugins: make(map[string]Plugin, 4),
	}

	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

// Select to choice some plugin to finish the task.
func (m *PluginManager) Select(ctx context.Context, query string) ([]PluginContext, error) {

	prompt := m.makePrompt(query)

	answer, err := m.chatWithLlm(ctx, prompt)
	if err != nil {
		logrus.Errorf("chat with llm error: %v", err)
		return nil, err
	}

	pluginCtxs := m.choicePlugins(answer)

	// for debug
	for _, c := range pluginCtxs {
		logrus.Debugf("query: %s choice plugins: %s input: %s", query, c.GetName(), c.Input)
	}

	return pluginCtxs, nil
}

func (m *PluginManager) makePrompt(query string) string {

	tools := m.makeTaskList()

	prompt := fmt.Sprintf(`You are an helpful and kind assistant to answer questions that can use tools to interact with real world and get access to the latest information.
	You will performs one task based on the following object:
	%s

	You can call one of the following functions:
	%s

	In each response, you must start with a function call like Tool name and args, split by ':',like:
	Google: query
	Weather:

	Don't explain why you use a tool. If you cannot figure out the answer, you say 'I don’t know'.

	Select only the corresponding tool and do not return any results.`,

		query,
		tools,
	)

	return prompt
}

func (m *PluginManager) makeTaskList() string {

	lines := make([]string, 0, len(m.plugins))

	for _, p := range m.plugins {

		line := fmt.Sprintf(
			`%s, Input Example: %s, It works as: %s`,
			p.GetName(),
			p.GetInputExample(),
			p.GetDesc(),
		)

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m *PluginManager) chatWithLlm(ctx context.Context, query string) (string, error) {
	messages := []llm.LlmMessage{
		{
			Role:    llm.RoleUser,
			Content: query,
		},
	}

	answer, err := m.llmer.Chat(ctx, messages)
	if err != nil {
		return "", err
	}

	// logrus.Debugf("query: %s\n  answer: %+v", query, answer)

	return answer.Content, nil
}

func (m *PluginManager) choicePlugins(answer string) []PluginContext {

	lines := strings.Split(answer, "\n")

	pluginContexts := make([]PluginContext, 0, len(lines))

	for _, line := range lines {
		logrus.Debugf("select one line: %s", line)

		if line == `I don’t know.` {
			continue
		}

		// Split by space
		// IF only ONE column, it's function name without args.
		// IF TWO column, it's function name with args.

		ss := strings.Split(line, ":")
		if len(ss) == 0 {
			logrus.Warnf("answer line invalid: %s", line)
			continue
		}

		name := strings.TrimSpace(strings.ToLower(ss[0]))
		var input string
		if len(ss) == 2 {
			input = strings.TrimSpace(ss[1])
		}

		if p, ok := m.plugins[name]; ok {

			logrus.Debugf("choice one plug with args: plugin=%v args=%v", name, input)

			pluginCtx := PluginContext{
				Plugin: p,
				Input:  input,
			}

			pluginContexts = append(pluginContexts, pluginCtx)
		}
	}

	return pluginContexts
}
