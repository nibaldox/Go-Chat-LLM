package api

import (
    "bufio"
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"

    "github.com/tuuser/go-ollama-tui/internal/model"
) // keep

// Client minimal Ollama HTTP client.
type Client struct {
    BaseURL string
    HTTP    *http.Client
}

func New(baseURL string) *Client {
    return &Client{BaseURL: baseURL, HTTP: &http.Client{}}
}

// Chat streams or gets a full response depending on req.Stream.
// Models returns available model names from Ollama.
func (c *Client) Models(ctx context.Context) ([]string, error) {
    r, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/tags", nil)
    if err != nil {
        return nil, err
    }
    resp, err := c.HTTP.Do(r)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var raw struct {
        Models []struct {
            Name string `json:"name"`
        } `json:"models"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
        return nil, err
    }
    names := make([]string, 0, len(raw.Models))
    for _, m := range raw.Models {
        names = append(names, m.Name)
    }
    return names, nil
}

func (c *Client) Chat(ctx context.Context, req model.ChatRequest) (<-chan model.ChatResponse, error) {
    b, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }
    r, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/chat", bytes.NewReader(b))
    if err != nil {
        return nil, err
    }
    r.Header.Set("Content-Type", "application/json")
    resp, err := c.HTTP.Do(r)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, io.ErrUnexpectedEOF
    }

    ch := make(chan model.ChatResponse)
    go func() {
        defer resp.Body.Close()
        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            var raw struct {
                Done    bool `json:"done"`
                Message struct {
                    Content string `json:"content"`
                } `json:"message"`
            }
            if err := json.Unmarshal(scanner.Bytes(), &raw); err == nil {
                ch <- model.ChatResponse{Done: raw.Done, Content: raw.Message.Content}
                if raw.Done {
                    break
                }
            }
        }
        close(ch)
    }()
    return ch, nil
}
