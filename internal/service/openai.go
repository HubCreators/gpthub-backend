package service

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
	model  string
}

func NewOpenAIService(client *openai.Client, model string) *OpenAIService {
	return &OpenAIService{
		client: client,
		model:  model,
	}
}

func (s *OpenAIService) Communicate(content string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}

	resp, err := s.client.CreateChatCompletion(context.Background(), req)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
