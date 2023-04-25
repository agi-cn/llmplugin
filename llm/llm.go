package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Role string

const (
	RoleUser      Role = openai.ChatMessageRoleUser
	RoleAssistant Role = openai.ChatMessageRoleAssistant
	RoleSystem    Role = openai.ChatMessageRoleSystem
)

func (r Role) String() string {
	return string(r)
}

type LlmMessage struct {
	Role    Role
	Content string
}

type LlmAnswer struct {
	Role    string
	Content string
}

type LLMer interface {
	Chat(ctx context.Context, messages []LlmMessage) (*LlmAnswer, error)
}
