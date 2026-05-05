# Gochestra

Gochestra is a small Go library that wraps multiple LLM providers behind a single `Agent` interface. It lets you swap providers and models quickly while keeping the same invocation shape.

## Features

- Unified `Agent` interface for chat-style calls.
- Providers: OpenAI, Anthropic, Gemini, Groq.
- Simple, stateless `Invoke` API.

## Install

```bash
go get github.com/alokxcode/gochestra
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"

	"github.com/alokxcode/gochestra/agent"
)

func main() {
	ag := agent.New(agent.OpenAi, "YOUR_API_KEY")

	cfg := &agent.ChatConfig{
		Model:        "gpt-4o-mini",
		SystemPrompt: "You are a helpful assistant.",
		MaxToken:     200,
	}

	prompt := agent.Prompt{Text: "Say hello in one sentence."}

	res, err := ag.Invoke(prompt, nil, cfg, context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Content.Openai.Choices[0].Message.Content)
}
```

## Providers

Use the `agent.Provider` constants when creating an agent:

- `agent.OpenAi`
- `agent.Anthropic`
- `agent.Gemini`
- `agent.Groq`

Each provider requires a valid API key. You pass it directly to `agent.New`.

## Configuration

`agent.ChatConfig` controls the request:

- `Model` - provider-specific model name
- `SystemPrompt` - system message sent before user content
- `MaxToken` - max output tokens

## Notes

- The library is intentionally minimal and stateless.
