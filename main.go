package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/corani/cuddly-octo-pancake/client"
)

const model = "openai/gpt-4.1"

const (
	roleSystem    = "system"
	roleAssistant = "assistant"
	roleUser      = "user"
	roleFunction  = "function"
	roleTool      = "tool"
)

func printAvailableModels(llm *client.Client, ctx context.Context) error {
	models, err := llm.ListModels(ctx)
	if err != nil {
		return err
	}
	log.Printf("Available models:")
	for _, model := range models {
		log.Printf("Model: %s, Tags: %v", model.ID, model.Tags)
	}
	return nil
}

func performChatCompletion(llm *client.Client, ctx context.Context, userMessage string) error {
	messages := []client.Message{
		{
			Role:    roleSystem,
			Content: "You are a helpful assistant.",
		},
		{
			Role:    roleUser,
			Content: userMessage,
		},
	}
	chat, err := llm.CreateChat(ctx, &client.ChatRequest{
		Messages: messages,
	})
	if err != nil {
		return err
	}
	for _, choice := range chat.Choices {
		messages = append(messages, choice.Message)
	}
	for _, msg := range messages {
		log.Printf("%s: %s", msg.Role, msg.Content)
	}
	return nil
}

func main() {
	var (
		showModels  = flag.Bool("models", false, "Print available models and exit")
		userMessage = flag.String("message", "What is the capital of France?", "User message for chat completion")
		helpFlag    = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage of cuddly-octo-pancake:")
		flag.PrintDefaults()
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	llm, err := client.NewClient(token, model, nil)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	if *showModels {
		if err := printAvailableModels(llm, ctx); err != nil {
			log.Fatalf("Error listing models: %v", err)
		}
		return
	}

	log.Printf("Using model %s", model)

	if err := performChatCompletion(llm, ctx, *userMessage); err != nil {
		log.Fatalf("Error creating chat: %v", err)
	}
}
