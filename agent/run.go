package agent

import (
	"strings"
	"context"

	"google.golang.org/genai"
)

const ModelName = "gemini-3.5-flash-lite"

// Callbacks is the seam between the agent core and any UI (the course's AgentCallbacks).
type Callbacks struct {
	OnToken    func(string)
	OnComplete func(string)
}

// Run does one streaming turn and returns the updated history.
// Lesson 4: single turn, no tools yet. The tool loop comes in lesson 14+.
func Run(ctx context.Context, client *genai.Client, userMessage string, history []*genai.Content, cb Callbacks) ([]*genai.Content, error) {
	// Step 1: save the user's message. Gemini is stateless, so history is its only memory.
	history = append(history, genai.NewContentFromText(userMessage, genai.RoleUser))

	// Step 2: send the FULL history and get the reply back in pieces (a Go iterator).
	// The loop body runs once per piece; we print each one and glue them into `full`.
	var full strings.Builder
	for chunk, err := range client.Models.GenerateContentStream(ctx, ModelName, history, nil) {
		if err != nil {
			return history, err
		}
		t := chunk.Text()
		full.WriteString(t)
		if cb.OnToken != nil {
			cb.OnToken(t)
		}
	}

	// Step 3: save the model's reply too. Skip this and Gemini forgets what it said
	// on the next turn, because we only ever send what is in history.
	history = append(history, genai.NewContentFromText(full.String(), genai.RoleModel))

	// Step 4: tell the UI we're done and hand the updated history back to the caller.
	if cb.OnComplete != nil {
		cb.OnComplete(full.String())
	}
	return history, nil
}
