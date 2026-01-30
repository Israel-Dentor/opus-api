package types

import "encoding/json"

// ClaudeRequest represents a Claude API request
type ClaudeRequest struct {
	Model       string                 `json:"model"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Messages    []ClaudeMessage        `json:"messages"`
	System      interface{}            `json:"system,omitempty"` // string or []ClaudeSystemMessage
	Tools       []ClaudeTool           `json:"tools,omitempty"`
	ToolChoice  interface{}            `json:"tool_choice,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	TopP        float64                `json:"top_p,omitempty"`
	TopK        int                    `json:"top_k,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeMessage represents a message in Claude API
type ClaudeMessage struct {
	Role    string      `json:"role"` // "user" or "assistant"
	Content interface{} `json:"content"` // string or []ClaudeContentBlock
}

// ClaudeSystemMessage represents a system message with cache control
type ClaudeSystemMessage struct {
	Type         string                 `json:"type"` // "text"
	Text         string                 `json:"text"`
	CacheControl map[string]interface{} `json:"cache_control,omitempty"`
}

// ClaudeContentBlock interface for different content types
type ClaudeContentBlock interface {
	GetType() string
}

// ClaudeContentBlockText represents text content
type ClaudeContentBlockText struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

func (c ClaudeContentBlockText) GetType() string { return c.Type }

// ClaudeContentBlockImage represents image content
type ClaudeContentBlockImage struct {
	Type   string `json:"type"` // "image"
	Source struct {
		Type      string `json:"type"` // "base64"
		MediaType string `json:"media_type"`
		Data      string `json:"data"`
	} `json:"source"`
}

func (c ClaudeContentBlockImage) GetType() string { return c.Type }

// ClaudeContentBlockToolUse represents tool use
type ClaudeContentBlockToolUse struct {
	Type  string                 `json:"type"` // "tool_use"
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

func (c ClaudeContentBlockToolUse) GetType() string { return c.Type }

// ClaudeContentBlockToolResult represents tool result
type ClaudeContentBlockToolResult struct {
	Type       string      `json:"type"` // "tool_result"
	ToolUseID  string      `json:"tool_use_id"`
	Content    interface{} `json:"content"` // string or []ClaudeContentBlock
	IsError    bool        `json:"is_error,omitempty"`
}

func (c ClaudeContentBlockToolResult) GetType() string { return c.Type }

// ClaudeTool represents a tool definition
type ClaudeTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// UnmarshalJSON custom unmarshaler for ClaudeMessage.Content
func (m *ClaudeMessage) UnmarshalJSON(data []byte) error {
	type Alias ClaudeMessage
	aux := &struct {
		Content json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(aux.Content, &str); err == nil {
		m.Content = str
		return nil
	}

	// Try to unmarshal as array of content blocks
	var blocks []json.RawMessage
	if err := json.Unmarshal(aux.Content, &blocks); err != nil {
		return err
	}

	var contentBlocks []ClaudeContentBlock
	for _, block := range blocks {
		var typeCheck struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(block, &typeCheck); err != nil {
			continue
		}

		switch typeCheck.Type {
		case "text":
			var textBlock ClaudeContentBlockText
			if err := json.Unmarshal(block, &textBlock); err == nil {
				contentBlocks = append(contentBlocks, textBlock)
			}
		case "image":
			var imageBlock ClaudeContentBlockImage
			if err := json.Unmarshal(block, &imageBlock); err == nil {
				contentBlocks = append(contentBlocks, imageBlock)
			}
		case "tool_use":
			var toolUseBlock ClaudeContentBlockToolUse
			if err := json.Unmarshal(block, &toolUseBlock); err == nil {
				contentBlocks = append(contentBlocks, toolUseBlock)
			}
		case "tool_result":
			var toolResultBlock ClaudeContentBlockToolResult
			if err := json.Unmarshal(block, &toolResultBlock); err == nil {
				contentBlocks = append(contentBlocks, toolResultBlock)
			}
		}
	}

	m.Content = contentBlocks
	return nil
}