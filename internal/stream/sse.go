package stream

import (
	"encoding/json"
	"fmt"
)

// FormatSSE formats an SSE event
func FormatSSE(event string, data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return fmt.Sprintf("event: %s\ndata: %s\n\n", event, string(jsonData))
}

// SSE event data structures

type MessageStartEvent struct {
	Type    string       `json:"type"`
	Message MessageStart `json:"message"`
}

type MessageStart struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Role         string                 `json:"role"`
	Content      []interface{}          `json:"content"`
	Model        string                 `json:"model"`
	StopReason   interface{}            `json:"stop_reason"`
	StopSequence interface{}            `json:"stop_sequence"`
	Usage        map[string]int         `json:"usage"`
}

type ContentBlockStartEvent struct {
	Type         string      `json:"type"`
	Index        int         `json:"index"`
	ContentBlock interface{} `json:"content_block"`
}

type TextContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ToolUseContentBlock struct {
	Type  string                 `json:"type"`
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

type ContentBlockDeltaEvent struct {
	Type  string      `json:"type"`
	Index int         `json:"index"`
	Delta interface{} `json:"delta"`
}

type TextDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type InputJSONDelta struct {
	Type        string `json:"type"`
	PartialJSON string `json:"partial_json"`
}

type ContentBlockStopEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type MessageDeltaEvent struct {
	Type  string                 `json:"type"`
	Delta map[string]interface{} `json:"delta"`
	Usage map[string]int         `json:"usage"`
}

type MessageStopEvent struct {
	Type string `json:"type"`
}