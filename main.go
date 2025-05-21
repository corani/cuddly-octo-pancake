package main

import (
	"context"
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

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	llm, err := client.NewClient(token, model, nil)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	models, err := llm.ListModels(ctx)
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}

	log.Printf("Available models:")

	for _, model := range models {
		log.Printf("Model: %s, Tags: %v", model.ID, model.Tags)
	}

	log.Printf("Using model %s", model)

	messages := []client.Message{
		{
			Role:    roleSystem,
			Content: "You are a helpful assistant.",
		},
		{
			Role:    roleUser,
			Content: "What is the capital of France?",
		},
	}
	chat, err := llm.CreateChat(ctx, &client.ChatRequest{
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Error creating chat: %v", err)
	}

	for _, choice := range chat.Choices {
		messages = append(messages, choice.Message)
	}

	for _, msg := range messages {
		log.Printf("%s: %s", msg.Role, msg.Content)
	}
}
