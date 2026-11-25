package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/FadyGamilM/build-ai-agent-from-scratch/agent"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	_ = godotenv.Load() // ponytail: .env optional, real env vars also work
	ctx := context.Background()

	// Reads GEMINI_API_KEY / GOOGLE_API_KEY from the environment.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{Backend: genai.BackendGeminiAPI})
	if err != nil {
		log.Fatal(err)
	}

	var history []*genai.Content
	in := bufio.NewScanner(os.Stdin)
	for fmt.Print("\n> "); in.Scan(); fmt.Print("\n> ") {
		history, err = agent.Run(ctx, client, in.Text(), history, agent.Callbacks{
			OnToken: func(s string) { fmt.Print(s) },
		})
		if err != nil {
			fmt.Println("\nerror:", err)
		}
	}
}
