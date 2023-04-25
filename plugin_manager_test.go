package llmplugin

import (
	"context"
	"testing"

	"github.com/agi-cn/llmplugin/llm"
	"github.com/agi-cn/llmplugin/llm/openai"
	"github.com/stretchr/testify/assert"
)

func TestManagerSelectPlugin(t *testing.T) {
	manager := newChatGPTManager()

	t.Run("Digital Computing", func(t *testing.T) {
		choice, err := manager.Select(context.Background(), "一加二等于几？")
		assert.NoError(t, err)

		assert.NotEmpty(t, choice)
		assert.True(t, includePlugin(choice, "Calculator"))
	})

	t.Run("Query Weather", func(t *testing.T) {
		choice, err := manager.Select(context.Background(), "今天天气如何？")
		assert.NoError(t, err)

		assert.NotEmpty(t, choice)
		assert.True(t, includePlugin(choice, "Weather"))
	})

}

func includePlugin(plugins []Plugin, target string) bool {
	for _, p := range plugins {
		if p.GetName() == target {
			return true
		}
	}

	return false
}

func TestChoicePlugins(t *testing.T) {

	plugins := newPlugins()
	manager := NewPluginManager(nil, WithPlugins(plugins))

	t.Run("Choice Calculator", func(t *testing.T) {

		answer := "Calculator \"1+2\""
		got := manager.choicePlugins(answer)

		assert.True(t,
			includePlugin(got, "Calculator"))
	})

	t.Run("Choice Weather", func(t *testing.T) {
		answer := "Weather"
		got := manager.choicePlugins(answer)

		assert.True(t,
			includePlugin(got, "Weather"))

	})
}

func newChatGPTManager() *PluginManager {
	var llmer llm.LLMer
	{
		token := "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"
		llmer = openai.NewChatGPT(token)
	}

	plugins := newPlugins()

	return NewPluginManager(llmer, WithPlugins(plugins))
}

func newPlugins() []Plugin {
	plugins := []Plugin{
		&SimplePlugin{
			Name: "Weather",
			Desc: "Can check the weather forecast",
			DoFunc: func(ctx context.Context, query string) (answer string, err error) {
				answer = "Call Weather Plugin"
				return
			},
		},
		&SimplePlugin{
			Name: "Calculator",
			Desc: "A calculator, capable of performing mathematical calculations, where the input is a description of a mathematical expression and the return is the result of the calculation. For example: the input is: one plus two, the return is three.",
			DoFunc: func(ctx context.Context, query string) (answer string, err error) {
				answer = "Call Calculator Plugin"
				return
			},
		},
	}
	return plugins
}
