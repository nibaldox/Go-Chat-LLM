package ui

import (
    "context"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/bubbles/spinner"
    "strings"
    "fmt"
    "github.com/charmbracelet/lipgloss"
    list "github.com/charmbracelet/bubbles/list"
    "time"
    "github.com/charmbracelet/glamour"
    "github.com/tuuser/go-ollama-tui/internal/mcp"

    "github.com/tuuser/go-ollama-tui/internal/api"
    "github.com/tuuser/go-ollama-tui/internal/model"
)

// ChatModel implements Bubble Tea model.
type ChatModel struct {
    ctx    context.Context
    client *api.Client

    input    textinput.Model
    viewport viewport.Model
    spinner  spinner.Model
    history  []string
    loading  bool
    width, height int
    modelName   string
    tokenCount  int
    currentAI string
    streamStart time.Time
    tools []mcp.Tool
    selecting bool
    modelList list.Model
    showHelp bool
    lastDuration time.Duration
    loadingFrame int
    isStreamingAI bool
    
}


func NewChatModel(ctx context.Context) ChatModel {
    client := api.New("http://localhost:11434")
    models, _ := client.Models(ctx)
    items := make([]list.Item, len(models))
    for i, m := range models {
        items[i] = listItem(m)
    }
    l := list.New(items, list.NewDefaultDelegate(), 30, 10)
    l.Title = "Selecciona modelo"
    l.SetShowStatusBar(false)
    l.SetFilteringEnabled(false)
    tools, _ := mcp.Load("docs/tools.json")
    ti := textinput.New()
    ti.Focus()
    sp := spinner.New()
    return ChatModel{
        ctx:    ctx,
        client:  api.New("http://localhost:11434"),
        input:   ti,
        viewport: viewport.New(80, 20),
        spinner:  sp,
        history:  []string{},
        selecting: len(items) > 0,
        modelList: l,
        modelName: "gemma3:4b-it-qat",
        currentAI: "",
        tools: tools,
    }
}

// Init initial command.
func (m ChatModel) Init() tea.Cmd {
    return textinput.Blink
}

// Update handles messages.
func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.viewport.Width = m.width - 2
        m.viewport.Height = m.height - 4 // status+input (no header)
        m.input.Width = m.width - 2
        
        // Reflow existing messages to new width
        m.viewport.SetContent(strings.Join(m.history, "\n"))
        return m, nil
    case tea.KeyMsg:
        if m.selecting {
            var cmd tea.Cmd
            m.modelList, cmd = m.modelList.Update(msg)
            switch msg.String() {
            case "enter":
                if sel, ok := m.modelList.SelectedItem().(listItem); ok {
                    m.modelName = string(sel)
                }
                m.selecting = false
            case "ctrl+c", "esc":
                m.selecting = false
            }
            return m, cmd
        }
        switch msg.String() {
        case "ctrl+c", "esc":
            return m, tea.Quit
        case "ctrl+h", "f1":
            m.showHelp = !m.showHelp
            return m, nil
        case "ctrl+m":
            m.selecting = true
            return m, nil
        case "ctrl+l":
            m.history = []string{}
            m.tokenCount = 0
            m.viewport.SetContent("")
            return m, nil
        case "enter":
            prompt := m.input.Value()
            m.input.Reset()
            
            // Responsive user message - max 80% of viewport width
            maxWidth := int(float64(m.viewport.Width) * 0.8)
            if maxWidth < 20 {
                maxWidth = 20
            }
            
            userMsg := userStyle.Width(maxWidth).Render("💬 You: "+prompt)
            aligned := lipgloss.NewStyle().Width(m.viewport.Width).Align(lipgloss.Right).Render(userMsg)
            m.history = append(m.history, aligned)
            
            // Responsive separator
            sepWidth := m.viewport.Width / 3 // Make separator 1/3 of viewport width
            if sepWidth < 3 {
                sepWidth = 3
            }
            separator := strings.Repeat("─", sepWidth)
            centeredSep := lipgloss.NewStyle().Width(m.viewport.Width).Align(lipgloss.Center).Render(separator)
            m.history = append(m.history, messageSeparatorStyle.Render(centeredSep))
            m.viewport.SetContent(strings.Join(m.history, "\n"))
            m.viewport.GotoBottom()
            m.loading = true
            m.tokenCount = 0
            m.loadingFrame = 0
            m.isStreamingAI = false
            m.streamStart = time.Now()
            return m, m.callOllama(prompt)
        }
    case ChunkMsg:
        chunk := msg.Chunk

        if chunk.Done {
            m.loading = false
            m.currentAI = ""
            m.isStreamingAI = false
            if !m.streamStart.IsZero() {
                m.lastDuration = time.Since(m.streamStart)
            }
            m.streamStart = time.Time{}
            return m, nil
        }
        m.loading = true
        m.loadingFrame++
        m.currentAI += chunk.Content
        
        // Responsive AI message - max 85% of viewport width
        maxWidth := int(float64(m.viewport.Width) * 0.85)
        if maxWidth < 20 {
            maxWidth = 20
        }
        
        rendered := renderMD(m.currentAI, maxWidth-8) // Account for padding and prefix
        aiPrefix := "🤖 AI:"
        aiMessage := aiStyle.Width(maxWidth).Render(aiPrefix + " " + rendered)
        
        // Use flag to track if we're streaming AI response
        if m.isStreamingAI {
            // Update existing AI message
            m.history[len(m.history)-1] = aiMessage
        } else {
            // Start new AI message
            m.history = append(m.history, aiMessage)
            m.isStreamingAI = true
        }
        m.viewport.SetContent(strings.Join(m.history, "\n"))
        m.viewport.GotoBottom()
        m.tokenCount += len(strings.Fields(chunk.Content))
        return m, streamNext(msg.Stream)
    }
    var cmd tea.Cmd
    m.input, cmd = m.input.Update(msg)
    m.viewport, _ = m.viewport.Update(msg)
    return m, cmd
}

type ChunkMsg struct {
    Chunk  model.ChatResponse
    Stream <-chan model.ChatResponse
}

func streamNext(ch <-chan model.ChatResponse) tea.Cmd {
    return func() tea.Msg {
        chunk, ok := <-ch
        if !ok {
            return ChunkMsg{Chunk: model.ChatResponse{Done: true}, Stream: nil}
        }
        return ChunkMsg{Chunk: chunk, Stream: ch}
    }
}

func (m ChatModel) callOllama(prompt string) tea.Cmd {
    return func() tea.Msg {
        stream, err := m.client.Chat(m.ctx, model.ChatRequest{
            Model:   m.modelName,
            Stream:  true,
            Messages: []model.ChatMessage{{Role: "user", Content: prompt}},
        })
        if err != nil {
            return model.ChatResponse{Done: true, Content: "[error]"}
        }
        return func() tea.Msg {
            chunk, ok := <-stream
            if !ok {
                return ChunkMsg{Chunk: model.ChatResponse{Done: true}, Stream: nil}
            }
            return ChunkMsg{Chunk: chunk, Stream: stream}
        }()
    }
}

// View renders.
func (m ChatModel) View() string {
    if m.selecting {
        return m.modelList.View()
    }
    
    // Body
    body := m.viewport.View()

    // Stats
    tps := 0.0
    if !m.streamStart.IsZero() {
        if sec := time.Since(m.streamStart).Seconds(); sec > 0 {
            tps = float64(m.tokenCount) / sec
        }
    }
    // Enhanced status bar with color coding
    var statusText string
    statusIcon := "📊"
    if m.loading {
        statusIcon = "⚡"
    }
    
    modelPart := infoStyle.Render(m.modelName)
    tokenPart := fmt.Sprintf("tokens: %s", infoStyle.Render(fmt.Sprintf("%d", m.tokenCount)))
    speedPart := fmt.Sprintf("%.1f t/s", tps)
    
    if tps > 50 {
        speedPart = successStyle.Render(speedPart)
    } else if tps > 20 {
        speedPart = infoStyle.Render(speedPart)
    } else {
        speedPart = helpStyle.Render(speedPart)
    }
    
    // Always show full stats in footer
    dur := m.lastDuration
    if m.loading && !m.streamStart.IsZero() {
        dur = time.Since(m.streamStart)
    }
    durPart := helpStyle.Render(dur.Truncate(time.Millisecond).String())
    statusText = fmt.Sprintf("%s %s • %s • %s • tiempo: %s", statusIcon, modelPart, tokenPart, speedPart, durPart)
    status := statusBarStyle.Width(m.width).Render(statusText)

    // Help panel
    helpPanel := ""
    if m.showHelp {
        helpContent := []string{
            "📋 Atajos de teclado:",
            "",
            shortcutStyle.Render("Enter") + "     → Enviar mensaje",
            shortcutStyle.Render("Ctrl+M") + "    → Cambiar modelo",
            shortcutStyle.Render("Ctrl+L") + "    → Limpiar chat",
            shortcutStyle.Render("Ctrl+H/F1") + " → Toggle ayuda",
            shortcutStyle.Render("↑/↓") + "      → Scroll mensajes",
            shortcutStyle.Render("Ctrl+C") + "   → Salir",
            "",
            helpStyle.Render("Presiona Ctrl+H para ocultar"),
        }
        helpText := strings.Join(helpContent, "\n")
        // Make help panel responsive
        maxHelpWidth := m.width - 4
        if maxHelpWidth < 40 {
            maxHelpWidth = 40
        }
        helpPanel = "\n" + errorStyle.Width(maxHelpWidth).Render(helpText)
    }

    // Footer
    footer := "\n" + m.input.View()

    return body + helpPanel + "\n" + status + footer
}

type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }

func renderMD(s string, width int) string {
    r, err := glamour.NewTermRenderer(
        glamour.WithWordWrap(width),
        glamour.WithAutoStyle(),
    )
    if err != nil {
        return s
    }
    out, err := r.Render(s)
    if err != nil {
        return s
    }
    return out
}
