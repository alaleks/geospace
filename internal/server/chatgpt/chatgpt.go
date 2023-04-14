// chatgpt package allows you to use api openai (language model chatGPT).
package chatgpt

import (
	"context"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	token       = "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" // access token
	temperature = 0.7                                                   // temperature for control the level of randomness
)

// ChatGPT contains client of openai.
type ChatGPT struct {
	client *openai.Client
}

// New is the constructor of ChatGPT.
func New() *ChatGPT {
	return &ChatGPT{
		client: openai.NewClient(token),
	}
}

// Use is for the possibility of using the language model ChatGPT.
func (c *ChatGPT) Use(prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: temperature,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	var result string

	if err != nil {
		return result, err
	}

	for _, msg := range resp.Choices {
		if len(strings.TrimSpace(msg.Message.Content)) > 0 {
			result += msg.Message.Content
		}
	}

	return result, nil
}
