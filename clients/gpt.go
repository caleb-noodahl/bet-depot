package clients

import (
	"github.com/caleb-noodahl/bet-depot/config"
	"context"

	"github.com/sashabaranov/go-openai"
)

type GPTClient struct {
	conf   *config.APIConf
	client *openai.Client
}

func NewGPTClient(conf *config.APIConf) (*GPTClient, error) {
	return &GPTClient{
		conf:   conf,
		client: openai.NewClient(conf.OpenApiKey),
	}, nil
}

func (g *GPTClient) Prompt(role, content string) (string, error) {
	resp, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4oMini,
			MaxTokens: 300,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: role,
				}, {
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	return resp.Choices[0].Message.Content, err
}
