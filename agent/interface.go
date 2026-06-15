package agent

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go/v3"
	openaiOption "github.com/openai/openai-go/v3/option"
	"google.golang.org/genai"
)

type ChatConfig struct {
	Model string
	SystemPrompt string
	MaxToken int
}

type Prompt struct {
	Text string
	Image string
	Clipboard string
}

type Message struct {
	Role string
	Content Content
}
type Content map[string]string
type ChatHistory []Message

type ResMessageType struct {
	Anthropic *anthropic.Message
	Gemini *genai.GenerateContentResponse
	Openai *openai.ChatCompletion
	Groq *openai.ChatCompletion
}
type Res struct {
	Content ResMessageType
	Error error
}
type Agent interface {
	Invoke(prompt Prompt,history ChatHistory, cfg *ChatConfig, ctx context.Context ) (*Res,error)
}
type Client struct {
	Gemini genai.Client
	Anthropic anthropic.Client
	Openai openai.Client
}


type Provider string
const (
	Gemini Provider = "gemini"
	Anthropic Provider = "anthropic"
	OpenAi Provider = "openai"
	Groq Provider = "groq"
)


func New(provider Provider, api_key string) Agent {
	switch provider {
	case Gemini:
		return initGemini(api_key)
	case Anthropic:
		return initAnthropic(api_key)
	case OpenAi:
		return initOpenai(api_key)
	case Groq:
		return initGroq(api_key)
	default:
		return &DefaultAdaptor{}
	}
}

// initialises gemini client
func initGemini(api_key string) Agent {

	cc := genai.ClientConfig{
		APIKey: api_key,
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &cc)
	if err != nil {
		fmt.Println("-----------------Error while init-------------", err)
	}

	adaptor := GeminiAdaptor{
		Client :  *client,
	}

	return &adaptor
}

// initialises anthropic client
func initAnthropic(api_key string) Agent {
	client := anthropic.NewClient(option.WithAPIKey(api_key))
	
	adaptor := AnthropicAdaptor{
		Client: client,
	}
	return &adaptor
}

// initialises openai Client
func initOpenai(api_key string) Agent {

	client := openai.NewClient(openaiOption.WithAPIKey(api_key))
	adaptor := OpenaiAdaptor {
		Client: client,
	}
	return &adaptor
}


// initialises groq 
func initGroq(api_key string) Agent {

	client := openai.NewClient(openaiOption.WithAPIKey(api_key), openaiOption.WithBaseURL("https://api.groq.com/openai/v1"))
	adaptor := OpenaiAdaptor {
		Client: client,
	}
	return &adaptor
}
