package google

import (
	"context"
	"os"
	"testing"

	"github.com/agi-cn/llmplugin/llm/openai"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGoogleWithoutSummary(t *testing.T) {

	_ = godotenv.Load() // ignore if file not exists

	var (
		apiToken = os.Getenv("GOOGLE_TOKEN")
		engineID = os.Getenv("GOOGLE_ENGINE_ID")
	)

	if apiToken == "" || engineID == "" {
		t.Skip("missing google env: GOOGLE_TOKEN or GOOGLE_ENGINE_ID. SKIP!")
	}

	g := NewGoogle(engineID, apiToken)

	answer, err := g.Do(context.Background(), "Who is Google Boss?")
	assert.NoError(t, err)

	assert.NotEmpty(t, answer)

	t.Logf("got answer: %v", answer)
}

func TestGoogleWithSummary(t *testing.T) {

	_ = godotenv.Load() // ignore if file not exists

	var (
		apiToken = os.Getenv("GOOGLE_TOKEN")
		engineID = os.Getenv("GOOGLE_ENGINE_ID")

		openaiToken = os.Getenv("OPENAI_TOKEN")
	)

	if apiToken == "" || engineID == "" || openaiToken == "" {
		t.Skip("missing google env: GOOGLE_TOKEN or GOOGLE_ENGINE_ID. SKIP!")
	}

	chatgpt := openai.NewChatGPT(openaiToken)
	g := NewGoogle(engineID, apiToken, WithSummarizer(chatgpt))

	answer, err := g.Do(context.Background(), "今天有什么nba比赛?")
	assert.NoError(t, err)

	assert.NotEmpty(t, answer)

	t.Logf("got answer: %v", answer)
}
