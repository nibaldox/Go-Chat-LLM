package model

import "github.com/tuuser/go-ollama-tui/internal/mcp"

// ChatMessage represents a single message in a chat.
type ChatMessage struct {
    Role    string `json:"role"` // "user" or "assistant"
    Content string `json:"content"`
}

// ChatRequest is the payload sent to Ollama.
type ChatRequest struct {
    Model   string        `json:"model"`
    Messages []ChatMessage `json:"messages"`
    Stream  bool          `json:"stream"`
    Tools   []mcp.Tool    `json:"tools,omitempty"`
}

// ChatResponse represents a single response chunk (stream=true) or full message.
type ChatResponse struct {
    Done    bool   `json:"done"`
    Content string `json:"message"`
}
