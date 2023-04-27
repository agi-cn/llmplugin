package google

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGoogle(t *testing.T) {

	_ = godotenv.Load() // ignore if file not exists

	var (
		apiToken = os.Getenv("GOOGLE_API_TOKEN")
		engineID = os.Getenv("GOOGLE_ENGINE_ID")
	)
	g := NewGoogle(engineID, apiToken)

	answer, err := g.Do(context.Background(), "Who is Google Boss?")
	assert.NoError(t, err)

	assert.NotEmpty(t, answer)

	t.Logf("got answer: %v", answer)
}
