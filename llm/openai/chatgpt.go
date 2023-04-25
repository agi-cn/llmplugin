package openai

import (
	"context"

	"github.com/agi-cn/llmplugin/llm"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	model  string
	client *openai.Client
}

type Option func(c *ChatGPT)

func WithModel(model string) Option {
	return func(c *ChatGPT) {
		c.model = model
	}
}

func NewChatGPT(token string, opts ...Option) *ChatGPT {

	client := openai.NewClient(token)

	chatgpt := &ChatGPT{
		model:  openai.GPT3Dot5Turbo,
		client: client,
	}

	for _, opt := range opts {
		opt(chatgpt)
	}

	return chatgpt
}

func (c ChatGPT) Chat(ctx context.Context, messages []llm.LlmMessage) (*llm.LlmAnswer, error) {

	chatGPTMessages := c.makeChatGPTMessage(messages)

	return c.send(ctx, chatGPTMessages)

}

func (c ChatGPT) makeChatGPTMessage(messages []llm.LlmMessage) []openai.ChatCompletionMessage {

	chatGPTMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, m := range messages {
		chatGPTMessages = append(chatGPTMessages, openai.ChatCompletionMessage{
			Role:    m.Role.String(),
			Content: m.Content,
		})
	}

	return chatGPTMessages
}

func (c ChatGPT) send(ctx context.Context, messages []openai.ChatCompletionMessage) (*llm.LlmAnswer, error) {

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	if choices := resp.Choices; len(choices) == 0 {
		return nil, errors.New("got empty ChatGPT response")
	}

	answer := c.convertLlmAnswer(resp)
	return answer, nil
}

func (c ChatGPT) convertLlmAnswer(openaiResp openai.ChatCompletionResponse) *llm.LlmAnswer {

	choices := openaiResp.Choices[0]

	return &llm.LlmAnswer{
		Role:    choices.Message.Role,
		Content: choices.Message.Content,
	}
}
