package stream

import "strings"

// ToolTagPrefixes are possible tool call tag prefixes
var ToolTagPrefixes = []string{
	"<function_calls",
	"<tool>",
	"<tool ",
	"<tools>",
	"<tools ",
}

// TextBuffer manages text buffering for streaming
type TextBuffer struct {
	PendingText      string
	ToolCallDetected bool
}

// NewTextBuffer creates a new text buffer
func NewTextBuffer() *TextBuffer {
	return &TextBuffer{
		PendingText:      "",
		ToolCallDetected: false,
	}
}

// Add adds text to the buffer
func (b *TextBuffer) Add(text string) {
	b.PendingText += text
}

// FlushSafeText flushes safe text that won't be part of a tool call tag
func (b *TextBuffer) FlushSafeText(emitFunc func(string)) {
	if b.PendingText == "" || b.ToolCallDetected {
		return
	}

	safeEndIndex := len(b.PendingText)

	// Check each tool call tag prefix
	for _, prefix := range ToolTagPrefixes {
		for i := 1; i <= len(prefix); i++ {
			partialTag := prefix[:i]
			idx := strings.LastIndex(b.PendingText, partialTag)
			if idx != -1 && idx+len(partialTag) == len(b.PendingText) {
				// Partial tag at end, keep it
				if idx < safeEndIndex {
					safeEndIndex = idx
				}
			}
		}
	}

	if safeEndIndex > 0 {
		safeText := b.PendingText[:safeEndIndex]
		if safeText != "" {
			emitFunc(safeText)
		}
		b.PendingText = b.PendingText[safeEndIndex:]
	}
}

// FlushAll flushes all pending text
func (b *TextBuffer) FlushAll(emitFunc func(string)) {
	if b.PendingText != "" {
		emitFunc(b.PendingText)
		b.PendingText = ""
	}
}

// Clear clears the buffer
func (b *TextBuffer) Clear() {
	b.PendingText = ""
}

// IsEmpty checks if buffer is empty
func (b *TextBuffer) IsEmpty() bool {
	return b.PendingText == ""
}