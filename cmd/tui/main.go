package main

import (
    "context"
    "log"

    tea "github.com/charmbracelet/bubbletea"

    "github.com/tuuser/go-ollama-tui/internal/ui"
)

func main() {
    ctx := context.Background()
    if err := tea.NewProgram(ui.NewChatModel(ctx)).Start(); err != nil {
        log.Fatal(err)
    }
}
