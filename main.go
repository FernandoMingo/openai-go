package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found or could not be loaded")
	}

	// Define flags
	model := flag.String("model", "gpt-3.5-turbo", "The model to use for the chat completion")
	temperature := flag.Float64("temperature", 0.5, "The temperature to use for the chat completion")
	maxTokens := flag.Int("max-tokens", 100, "The maximum number of tokens to use for the chat completion")
	apiKey := flag.String("api-key", "", "The API key to use for the chat completion")
	flag.Parse()

	// Get the API key from the flags. If empty, use the environment variable
	key := *apiKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}

	// The remaining arguments after the flags are the prompt
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Please provide a message to send to the chat completion \n Example usage: go run main.go [--model=model] [--temperature=val] \"your prompt here\"")
		os.Exit(1)
	}
	prompt := strings.Join(args, " ")

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       *model,
			Temperature: float32(*temperature),
			MaxTokens:   *maxTokens,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
