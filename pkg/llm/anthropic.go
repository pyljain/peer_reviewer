package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	anthropicCompletionsUrl = "https://api.anthropic.com/v1/messages"
)

type Anthropic struct {
	apiKey string
}

func NewAnthropic() (*Anthropic, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing Anthropic API key")
	}
	return &Anthropic{
		apiKey: apiKey,
	}, nil
}

func (a *Anthropic) GetCompletions(
	ctx context.Context,
	model string,
	temperature float64,
	messages []Message,
	systemPrompt string) (*Message, error) {
	body := AnthropicRequestBody{
		Model:       model,
		Temperature: temperature,
		Messages:    messages,
		MaxTokens:   1024,
		Stream:      false,
		System:      systemPrompt,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, anthropicCompletionsUrl, buf)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Add("x-api-key", a.apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	req.Header.Add("content-type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ar := AnthropicResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ar)
	if err != nil {
		return nil, err
	}

	if len(ar.Content) == 0 {
		return &Message{
			Role:    ar.Role,
			Content: "No comments",
		}, nil
	}

	return &Message{
		Role:    ar.Role,
		Content: ar.Content[0].Text,
	}, nil

}

type AnthropicRequestBody struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	System      string    `json:"system"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
	Messages    []Message `json:"messages"`
}

type AnthropicResponse struct {
	Content []AnthropicResponseContent `json:"content"`
	Model   string                     `json:"model"`
	Role    string                     `json:"role"`
	Usage   AnthropicResponseUsage     `json:"usage"`
}

type AnthropicResponseContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type AnthropicResponseUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
