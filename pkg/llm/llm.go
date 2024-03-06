package llm

import (
	"context"
	"fmt"
)

type LLM interface {
	GetCompletions(ctx context.Context, model string, temperature float64, messages []Message, systemPrompt string) (*Message, error)
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func GetLLM(name string) (LLM, error) {
	switch name {
	case "anthropic":
		return NewAnthropic()
	default:
		return nil, fmt.Errorf("LLM not supported")
	}
}
