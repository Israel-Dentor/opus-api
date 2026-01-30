package stream

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestTransformMorphToClaudeStream_IncompleteToolCall tests the case where
// upstream returns incomplete tool call XML
func TestTransformMorphToClaudeStream_IncompleteToolCall(t *testing.T) {
	// Read test data
	testData, err := os.ReadFile("../../test/fixtures/incomplete_tool_call.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Create input reader
	input := bytes.NewReader(testData)

	// Create output buffer
	var output bytes.Buffer

	// Transform
	err = TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Errorf("Transform failed: %v", err)
	}

	// Verify output
	outputStr := output.String()

	// Check for required events
	if !strings.Contains(outputStr, "event: message_start") {
		t.Error("Missing message_start event")
	}

	if !strings.Contains(outputStr, "event: message_stop") {
		t.Error("Missing message_stop event - this is the bug!")
	}

	// Check for content
	if !strings.Contains(outputStr, "I'll analyze this") {
		t.Error("Missing expected content")
	}

	// The incomplete tool call XML should be discarded
	// This is expected behavior - incomplete tool calls are not output
	if strings.Contains(outputStr, "function_calls") {
		t.Error("Incomplete tool call XML should not be in output")
	}

	// Count events
	eventCount := strings.Count(outputStr, "event:")
	t.Logf("Total events: %d", eventCount)

	// Print output for debugging
	t.Logf("Output:\n%s", outputStr)
}

// TestTransformMorphToClaudeStream_CompleteToolCall tests normal tool call
func TestTransformMorphToClaudeStream_CompleteToolCall(t *testing.T) {
	// This will be added later when we have a complete tool call example
	t.Skip("TODO: Add complete tool call test case")
}

// TestTransformMorphToClaudeStream_NoToolCall tests normal text response
func TestTransformMorphToClaudeStream_NoToolCall(t *testing.T) {
	testData := `data: {"type":"start"}

data: {"type":"start-step"}

data: {"type":"text-start","id":"0"}

data: {"type":"text-delta","id":"0","delta":"Hello"}

data: {"type":"text-delta","id":"0","delta":" world"}

data: {"type":"text-end","id":"0"}

data: {"type":"finish-step"}

data: {"type":"finish","finishReason":"stop"}

data: [DONE]

`

	input := strings.NewReader(testData)
	var output bytes.Buffer

	err := TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	outputStr := output.String()

	// Verify all required events
	requiredEvents := []string{
		"event: message_start",
		"event: content_block_start",
		"event: content_block_delta",
		"event: content_block_stop",
		"event: message_delta",
		"event: message_stop",
	}

	for _, event := range requiredEvents {
		if !strings.Contains(outputStr, event) {
			t.Errorf("Missing required event: %s", event)
		}
	}

	// Check content - note: content may be split across multiple deltas
	if !strings.Contains(outputStr, "Hello") || !strings.Contains(outputStr, "world") {
		t.Error("Missing expected content")
	}

	t.Logf("Output:\n%s", outputStr)
}