package main

import (
	"context"
	"fmt"

	"github.com/alokxcode/gochestra/agent"
)

func main() {
	Agent := agent.New("groq", "your_api_key")
	cfg := agent.ChatConfig{
		Model:        "llama-3.3-70b-versatile",
		MaxToken:     200,
		SystemPrompt: "Hello",
	}
	prmpt := agent.Prompt{
		Text: "Hii",
	}
	content := map[string]string{"text": "hello"}
	history := agent.Message{
		Role:    "user",
		Content: content,
	}
	memory := []agent.Message{history}
	res, err := Agent.Invoke(prmpt, memory, &cfg, context.Background())
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Println(res.Content.Openai.Choices[0].Message.Content)
}
