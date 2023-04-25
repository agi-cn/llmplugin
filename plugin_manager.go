package llmplugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/agi-cn/llmplugin/llm"

	"github.com/sirupsen/logrus"
)

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
func (m *PluginManager) Select(ctx context.Context, query string) ([]Plugin, error) {

	prompt := m.makePrompt(query)

	answer, err := m.chatWithLlm(ctx, prompt)
	if err != nil {
		logrus.Errorf("chat with llm error: %v", err)
		return nil, err
	}

	plugins := m.choicePlugins(answer)
	return plugins, nil
}

func (m *PluginManager) makePrompt(query string) string {

	prompt := fmt.Sprintf(`You are an helpful and kind assistant to answer questions that can use tools to interact with real world and get access to the latest information.
	You will performs one task based on the following object:
	%s

	You can call one of the following functions:

	- Calculator, INPUT: (expr string): ACT ON A calculator, capable of performing mathematical calculations, where the input is a description of a mathematical expression and the return is the result of the calculation. For example: the input is: one plus two, the return is three.
	- Weather, INPUT: no input: ACT ON You can check the weather forecast.

	In each response, you must start with a function call like Tool name and args, split by space,like:
	Google "query"
	Weather

	Don't explain why you use a tool. If you cannot figure out the answer, you say ’I don’t know’.

	You can choose one tool or multiple tools to accomplish the corresponding goal.
	Select only the corresponding tool and do not return any results.`,
		query,
	)

	return prompt
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

	logrus.Debugf("query: %s\n  answer: %+v", query, answer)

	return answer.Content, nil
}

func (m *PluginManager) choicePlugins(answer string) []Plugin {

	lines := strings.Split(answer, "\n")

	plugins := make([]Plugin, 0, len(lines))

	for _, line := range lines {

		// Split by space
		// IF only ONE column, it's function name without args.
		// IF TWO column, it's function name with args.

		ss := strings.Split(line, " ")
		if len(ss) == 0 {
			continue
		}

		n := strings.ToLower(ss[0])

		if p, ok := m.plugins[n]; ok {
			plugins = append(plugins, p)
		}
	}

	return plugins
}
