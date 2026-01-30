package types

// MorphRequest represents a MorphLLM API request
type MorphRequest struct {
	SandboxID string         `json:"sandboxId"`
	RepoRoot  string         `json:"repoRoot"`
	ID        string         `json:"id"`
	Messages  []MorphMessage `json:"messages"`
	Trigger   string         `json:"trigger"`
}

// MorphMessage represents a message in MorphLLM API
type MorphMessage struct {
	Parts []MorphPart `json:"parts"`
	ID    string      `json:"id"`
	Role  string      `json:"role"`  // "user" or "assistant"
	State string      `json:"state"` // "done" or "pending"
}

// MorphPart represents a part of a message
type MorphPart struct {
	Type  string `json:"type"` // "text"
	Text  string `json:"text"`
	State string `json:"state"` // "done" or "pending"
}
