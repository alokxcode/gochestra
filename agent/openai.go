package agent

import (
	"context"

	"github.com/openai/openai-go/v3"
)

type OpenaiAdaptor struct {
	Client openai.Client
}

func (a *OpenaiAdaptor) Invoke(prompt Prompt, history ChatHistory, cfg *ChatConfig, ctx context.Context) (*Res, error) {
	conversation := append(history, Message{
		Role: "user",
		Content: map[string]string{
			"text":      prompt.Text,
			"image":     prompt.Image,
			"clipboard": prompt.Clipboard,
		},
	})

	var messages []openai.ChatCompletionMessageParamUnion

	// system prompt goes as first message
	messages = append(messages, openai.SystemMessage(cfg.SystemPrompt))

	for _, r := range conversation {
		var text string
		for _, m := range r.Content {
			text += m
		}
		switch r.Role {
		case "user":
			messages = append(messages, openai.UserMessage(text))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(text))
		}
	}

	response, err := a.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     cfg.Model,
		Messages:  messages,
		MaxTokens: openai.Int(int64(cfg.MaxToken)),
	})
	if err != nil {
		return nil, err
	}

	return &Res{
		Content: ResMessageType{Openai: response},
	}, nil
}
