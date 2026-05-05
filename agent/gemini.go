package agent

import (
	"context"

	"google.golang.org/genai"
)

type GeminiAdaptor struct {
	Client genai.Client
}

func (a *GeminiAdaptor) Invoke(prompt Prompt, history ChatHistory, cfg *ChatConfig, ctx context.Context) (*Res, error) {
	conversation := append(history, Message{
		Role: "user",
		Content: map[string]string{
			"text":      prompt.Text,
			"image":     prompt.Image,
			"clipboard": prompt.Clipboard,
		},
	})

	// Build history slice (everything except the last user message)
	var chatHistory []*genai.Content
	for _, r := range conversation[:len(conversation)-1] {
		var parts []*genai.Part
		for _, m := range r.Content {
			if m == "" {
				continue
			}
			parts = append(parts, &genai.Part{Text: m})
		}
		role := r.Role
		if role == "assistant" {
			role = "model"
		}
		chatHistory = append(chatHistory, &genai.Content{
			Role:  role,
			Parts: parts,
		})
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: cfg.SystemPrompt}},
		},
		MaxOutputTokens: int32(cfg.MaxToken),
	}

	chat, err := a.Client.Chats.Create(ctx, string(cfg.Model), config, chatHistory)
	if err != nil {
		return nil, err
	}

	// Build final user turn parts
	lastMsg := conversation[len(conversation)-1]
	var userParts []genai.Part
	for _, m := range lastMsg.Content {
		if m == "" {
			continue
		}
		userParts = append(userParts, genai.Part{Text: m})
	}

	response, err := chat.SendMessage(ctx, userParts...)
	if err != nil {
		return nil, err
	}

	return &Res{
		Content: ResMessageType{Gemini: response},
	}, nil
}
