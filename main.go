package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type ChatCompleter interface {
	CreateChatCompletion(model string, temperature float64, maxTokens int, prompt string) (string, error)
}

type OpenAIClient struct {
	client *openai.Client
}

func (o *OpenAIClient) CreateChatCompletion(model string, temperature float64, maxTokens int, prompt string) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Temperature: float32(temperature),
			MaxTokens:   maxTokens,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func RunChatCompletion(completer ChatCompleter, model string, temperature float64, maxTokens int, prompt string) (string, error) {
	return completer.CreateChatCompletion(model, temperature, maxTokens, prompt)
}

func main() {

	model, temperature, maxTokens, apiKey, prompt, err := ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := &OpenAIClient{client: openai.NewClient(apiKey)}
	response, err := RunChatCompletion(client, model, temperature, maxTokens, prompt)
	if err != nil {
		log.Fatalf("ChatCompletion error: %v", err)
	}
	fmt.Println(response)
}

// BuildPrompt joins CLI args into a single prompt string.
func BuildPrompt(args []string) string {
	return strings.Join(args, " ")
}

// ValidateAPIKey checks if the API key is non-empty.
func ValidateAPIKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key is required")
	}
	return nil
}

func ParseConfig() (model string, temperature float64, maxTokens int, apiKey string, prompt string, err error) {
	modelFlag := flag.String("model", "gpt-3.5-turbo", "The model to use for the chat completion")
	temperatureFlag := flag.Float64("temperature", 0.5, "The temperature to use for the chat completion")
	maxTokensFlag := flag.Int("max-tokens", 100, "The maximum number of tokens to use for the chat completion")
	apiKeyFlag := flag.String("api-key", "", "The API key to use for the chat completion")
	flag.Parse()

	apiKey = *apiKeyFlag
	if apiKey == "" {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found or could not be loaded")
		}
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		err = fmt.Errorf("no API key provided. set --api-key or OPENAI_API_KEY in .env or environment")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		err = fmt.Errorf("please provide a message to send to the chat completion.\n example usage: go run main.go [--model=model] [--temperature=val] \"your prompt here\"")
		return
	}
	prompt = BuildPrompt(args)
	model = *modelFlag
	temperature = *temperatureFlag
	maxTokens = *maxTokensFlag
	return
}
