package agent

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
)

type AnthropicAdaptor struct {
	Client anthropic.Client
}

func (a *AnthropicAdaptor) Invoke(prompt Prompt, history ChatHistory, cfg *ChatConfig,ctx context.Context) (*Res,error) {

	content := map[string]string{"text":prompt.Text,"image":prompt.Image,"clipboard":prompt.Clipboard}
	
	conversation := append(history,Message{Role: "user",Content: content})
	var messages []anthropic.MessageParam
	// loops over each msgs - [ [ user : "hii" ], [ ai : "hello" ] ]
	for _, r := range conversation{
		var contentBlock []anthropic.ContentBlockParamUnion
		// loops over each different types of msgs - [ " user : 'what's the promlem in this code' ", IMG, CODE_SNIPPET ]
		for _, m := range r.Content {
			contentBlock = append(contentBlock, anthropic.NewTextBlock(m))
		}
		msg := anthropic.MessageParam{
			Content: contentBlock,
			Role:    anthropic.MessageParamRole(r.Role),
		}
		messages = append(messages, msg)

	}

	systemPromptStruct := anthropic.TextBlockParam{
		Text: cfg.SystemPrompt,
	}

	systemPromptSlice := []anthropic.TextBlockParam{systemPromptStruct}

	message := anthropic.MessageNewParams{
		Model:     cfg.Model,
		Messages:  messages,
		System:    systemPromptSlice,
		MaxTokens: int64(cfg.MaxToken),
	}
	response,err := a.Client.Messages.New(ctx, message)
	if err != nil {
		return nil,err
	}
	res := &Res{
		Content: ResMessageType{Anthropic: response},
		Error: err,
	}
	return res,nil
}
