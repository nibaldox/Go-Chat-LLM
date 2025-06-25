package mcp

import (
    "encoding/json"
    "io/ioutil"
)

// Tool definition per MCP spec.
type Tool struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Parameters  map[string]interface{} `json:"parameters"`
}

// Load tools from JSON file.
func Load(path string) ([]Tool, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var tools []Tool
    if err := json.Unmarshal(data, &tools); err != nil {
        return nil, err
    }
    return tools, nil
}
