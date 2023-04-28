package llmplugin

import (
	"context"
	"os"
	"testing"

	"github.com/agi-cn/llmplugin/llm"
	"github.com/agi-cn/llmplugin/llm/openai"
	"github.com/agi-cn/llmplugin/plugins/calculator"
	"github.com/agi-cn/llmplugin/plugins/google"
	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagerSelectPlugin(t *testing.T) {
	manager := newChatGPTManager()

	t.Run("Digital Computing", func(t *testing.T) {
		pluginCtxs, err := manager.Select(context.Background(), "10 add 20 equals ?")
		require.NoError(t, err)

		require.Equal(t, 1, len(pluginCtxs))
		require.True(t, includePlugin(pluginCtxs, "Calculator"))

		choice := pluginCtxs[0]

		answer, err := choice.Plugin.Do(context.Background(), choice.Input)
		require.NoError(t, err)

		assert.Equal(t, "30", answer)
	})

	t.Run("Query Weather", func(t *testing.T) {
		choice, err := manager.Select(context.Background(), "How is the weather today?")
		assert.NoError(t, err)

		assert.NotEmpty(t, choice)
		assert.True(t, includePlugin(choice, "Weather"))
	})

}

func includePlugin(pluginCtxs []PluginContext, target string) bool {
	for _, p := range pluginCtxs {
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

		answer := "Calculator: 1+4"
		got := manager.choicePlugins(answer)

		assert.True(t,
			includePlugin(got, "Calculator"))
	})

	t.Run("Choice Weather", func(t *testing.T) {
		answer := "Weather: "
		got := manager.choicePlugins(answer)

		assert.True(t,
			includePlugin(got, "Weather"))

	})

	t.Run("Choice Google", func(t *testing.T) {
		answer := `Google: 今天NBA比赛赛程表`
		got := manager.choicePlugins(answer)

		assert.True(t,
			includePlugin(got, "Google"))
	})
}

func newChatGPTManager() *PluginManager {
	_ = godotenv.Load() // ignore if file not exists

	var llmer llm.LLMer
	{
		token := os.Getenv("OPENAI_TOKEN")
		if len(token) == 0 {
			panic("empty openai token: set os env: OPENAI_TOKEN")
		}
		llmer = openai.NewChatGPT(token)
	}

	plugins := newPlugins()

	return NewPluginManager(llmer, WithPlugins(plugins))
}

func newPlugins() []Plugin {

	var (
		googleEngineID = os.Getenv("GOOGLE_ENGINE_ID")
		googleToken    = os.Getenv("GOOGLE_TOKEN")
	)

	plugins := []Plugin{
		&SimplePlugin{
			Name:         "Weather",
			InputExample: ``,
			Desc:         "Can check the weather forecast",
			DoFunc: func(ctx context.Context, query string) (answer string, err error) {
				answer = "Call Weather Plugin"
				return
			},
		},

		calculator.NewCalculator(),

		google.NewGoogle(googleEngineID, googleToken),
	}
	return plugins
}
