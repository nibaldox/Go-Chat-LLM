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
        m.viewport.Height = m.height - 6 // header+status+input
        m.input.Width = m.width - 2
        m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.history, "\n")))
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
            userMsg := userStyle.Render("💬 You: "+prompt)
            aligned := lipgloss.NewStyle().Width(m.viewport.Width).Align(lipgloss.Right).Render(userMsg)
            m.history = append(m.history, aligned)
            m.history = append(m.history, separatorStyle.Render(""))
            m.viewport.SetContent(strings.Join(m.history, "\n"))
            m.viewport.GotoBottom()
            m.loading = true
            m.tokenCount = 0
            m.streamStart = time.Now()
            return m, m.callOllama(prompt)
        }
    case ChunkMsg:
        chunk := msg.Chunk

        if chunk.Done {
            m.loading = false
            m.currentAI = ""
            m.streamStart = time.Time{}
            return m, nil
        }
        m.loading = true
        m.currentAI += chunk.Content
        rendered := renderMD(m.currentAI, m.viewport.Width-4)
        aiPrefix := "🤖 AI:"
        if len(m.history) > 0 && strings.Contains(m.history[len(m.history)-1], aiPrefix) {
            m.history[len(m.history)-1] = aiStyle.Render(aiPrefix + " " + rendered)
        } else {
            m.history = append(m.history, aiStyle.Render(aiPrefix+" "+rendered))
        }
        m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.history, "\n")))
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
    
    // Enhanced header with better styling
    header := headerStyle.Render(fmt.Sprintf("🚀 Chat LLM • %s", m.modelName))
    
    // Main chat body
    body := m.viewport.View()
    
    // Loading indicator with better styling
    if m.loading {
        loadingMsg := loadingStyle.Render("⏳ Generando respuesta...")
        body += "\n" + loadingMsg
    }
    
    // Enhanced status bar with better token info
    tps := 0.0
    if !m.streamStart.IsZero() {
        elapsed := time.Since(m.streamStart).Seconds()
        if elapsed > 0 {
            tps = float64(m.tokenCount) / elapsed
        }
    }
    
    // Status with better formatting and shortcuts
    statusContent := fmt.Sprintf("📊 %s • %s tokens • %s t/s", 
        infoStyle.Render(m.modelName),
        infoStyle.Render(fmt.Sprintf("%d", m.tokenCount)),
        infoStyle.Render(fmt.Sprintf("%.1f", tps)))
    
    shortcuts := helpStyle.Render("Ctrl+M modelo • Ctrl+L limpiar • Ctrl+C salir • ↑↓ scroll")
    
    status := statusStyle.Width(m.width).Render(statusContent + " • " + shortcuts)
    
    // Input with better styling
    footer := "\n" + m.input.View()
    
    return header + "\n" + body + "\n" + status + footer
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
