package stream

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestNativeToolCallConversion tests conversion of native MorphLLM tool calls
// to Claude API format
func TestNativeToolCallConversion(t *testing.T) {
	// Read test data - this contains native tool calls from MorphLLM
	testData, err := os.ReadFile("/Users/leokun/Desktop/opus-api/一个完整的任务日志/4_upstream_response.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	input := bytes.NewReader(testData)
	var output bytes.Buffer

	err = TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Errorf("Transform failed: %v", err)
	}

	outputStr := output.String()

	// Verify required events
	if !strings.Contains(outputStr, "event: message_start") {
		t.Error("Missing message_start event")
	}

	if !strings.Contains(outputStr, "event: message_stop") {
		t.Error("Missing message_stop event")
	}

	// Check for tool use events
	if !strings.Contains(outputStr, "event: content_block_start") {
		t.Error("Missing content_block_start event")
	}

	// Verify tool use content blocks
	if !strings.Contains(outputStr, `"type":"tool_use"`) {
		t.Error("Missing tool_use content blocks")
	}

	// warp_grep is a native MorphLLM tool and should be ignored
	// We should check for the XML tool calls (Glob, Bash) instead

	// Check for Glob tool calls
	if !strings.Contains(outputStr, `"name":"Glob"`) {
		t.Error("Missing Glob tool call - XML tools not converted")
	}

	// Check for Bash tool calls
	if !strings.Contains(outputStr, `"name":"Bash"`) {
		t.Error("Missing Bash tool call - XML tools not converted")
	}

	// Verify tool input is properly formatted as JSON
	if !strings.Contains(outputStr, `"input":{`) {
		t.Error("Tool input not properly formatted")
	}

	// Should NOT contain XML tool call tags in output
	if strings.Contains(outputStr, "function_calls") {
		t.Error("Output should not contain XML function_calls tags")
	}

	if strings.Contains(outputStr, "<invoke") {
		t.Error("Output should not contain XML invoke tags")
	}

	// Verify stop reason is tool_use
	if !strings.Contains(outputStr, `"stop_reason":"tool_use"`) {
		t.Error("Stop reason should be tool_use")
	}

	t.Logf("Total output length: %d bytes", len(outputStr))

	// Print first 2000 chars for inspection
	if len(outputStr) > 2000 {
		t.Logf("Output preview:\n%s\n...", outputStr[:2000])
	} else {
		t.Logf("Output:\n%s", outputStr)
	}
}

// TestMultipleFunctionCallsBlocks tests handling of multiple separate
// function_calls blocks (which is invalid but may occur)
func TestMultipleFunctionCallsBlocks(t *testing.T) {
	testData, err := os.ReadFile("../../test/fixtures/multiple_function_calls.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	input := bytes.NewReader(testData)
	var output bytes.Buffer

	err = TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Errorf("Transform failed: %v", err)
	}

	outputStr := output.String()

	// Should have proper event structure
	if !strings.Contains(outputStr, "event: message_start") {
		t.Error("Missing message_start event")
	}

	if !strings.Contains(outputStr, "event: message_stop") {
		t.Error("Missing message_stop event")
	}

	// Should convert XML tool calls to proper tool_use blocks
	toolUseCount := strings.Count(outputStr, `"type":"tool_use"`)
	t.Logf("Found %d tool_use blocks", toolUseCount)

	if toolUseCount == 0 {
		t.Error("Should have converted XML tool calls to tool_use blocks")
	}

	// Should NOT have XML in output
	if strings.Contains(outputStr, "function_calls") {
		t.Error("Output should not contain XML function_calls tags")
	}

	t.Logf("Output preview:\n%s", outputStr[:min(2000, len(outputStr))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestRealWarpGrepToolCall tests conversion of real warp_grep tool call
// from logs/1 to ensure parameters are properly preserved
func TestRealWarpGrepToolCall(t *testing.T) {
	// This test uses the actual SSE response from logs/1/morph_response.txt
	// which contains a warp_grep tool call with proper parameters
	testData, err := os.ReadFile("../../test/fixtures/real_warp_grep_call.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	input := bytes.NewReader(testData)
	var output bytes.Buffer

	err = TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Errorf("Transform failed: %v", err)
	}

	outputStr := output.String()

	// Verify basic structure
	if !strings.Contains(outputStr, "event: message_start") {
		t.Error("Missing message_start event")
	}

	if !strings.Contains(outputStr, "event: message_stop") {
		t.Error("Missing message_stop event")
	}

	// Check for tool use blocks
	if !strings.Contains(outputStr, `"type":"tool_use"`) {
		t.Error("Missing tool_use content blocks")
	}

	// CRITICAL: Verify that tool parameters are NOT empty
	// The bug is that input becomes {} instead of containing actual parameters
	if strings.Contains(outputStr, `"input":{}`) {
		t.Error("CRITICAL BUG: Tool input is empty {}! Parameters were lost during conversion")
	}

	// Verify Glob tool calls have pattern parameter
	if strings.Contains(outputStr, `"name":"Glob"`) {
		// Check that Glob calls have pattern parameter
		if !strings.Contains(outputStr, `"pattern"`) {
			t.Error("Glob tool call missing pattern parameter")
		}
		// Specific patterns from the real data
		if !strings.Contains(outputStr, "package.json") {
			t.Error("Missing expected pattern 'package.json' in Glob call")
		}
	}

	// Verify Bash tool calls have command parameter
	if strings.Contains(outputStr, `"name":"Bash"`) {
		if !strings.Contains(outputStr, `"command"`) {
			t.Error("Bash tool call missing command parameter")
		}
		// Check for specific commands from real data
		if !strings.Contains(outputStr, "find") || !strings.Contains(outputStr, "ls") {
			t.Error("Missing expected commands in Bash calls")
		}
	}

	// Verify stop reason
	if !strings.Contains(outputStr, `"stop_reason":"tool_use"`) {
		t.Error("Stop reason should be tool_use")
	}

	// Should NOT contain XML in output
	if strings.Contains(outputStr, "function_calls") {
		t.Error("Output should not contain XML function_calls tags")
	}

	t.Logf("Total output length: %d bytes", len(outputStr))

	// Print sample of output for debugging
	lines := strings.Split(outputStr, "\n")
	t.Logf("Total lines: %d", len(lines))

	// Find and print tool_use blocks
	for i, line := range lines {
		if strings.Contains(line, `"type":"tool_use"`) ||
			strings.Contains(line, `"input"`) ||
			strings.Contains(line, `"name":"Glob"`) ||
			strings.Contains(line, `"name":"Bash"`) {
			t.Logf("Line %d: %s", i, line)
		}
	}
}

// TestLatestMorphResponse tests the latest morph response to verify tool conversion
func TestLatestMorphResponse(t *testing.T) {
	testData, err := os.ReadFile("../../test/fixtures/latest_morph_response.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	input := bytes.NewReader(testData)
	var output bytes.Buffer

	err = TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Errorf("Transform failed: %v", err)
	}

	outputStr := output.String()

	// Check for tool use blocks
	toolUseCount := strings.Count(outputStr, `"type":"tool_use"`)
	t.Logf("Found %d tool_use blocks", toolUseCount)

	if toolUseCount == 0 {
		t.Error("Expected tool_use blocks but found none")
	}

	// Check for specific tools
	if !strings.Contains(outputStr, `"name":"Read"`) {
		t.Error("Missing Read tool call")
	}

	if !strings.Contains(outputStr, `"name":"Glob"`) {
		t.Error("Missing Glob tool call")
	}

	// Verify stop reason
	if !strings.Contains(outputStr, `"stop_reason":"tool_use"`) {
		t.Error("Stop reason should be tool_use, not end_turn")
	}

	// Print tool blocks for inspection
	lines := strings.Split(outputStr, "\n")
	t.Logf("Total lines: %d", len(lines))
	for i, line := range lines {
		if strings.Contains(line, `"type":"tool_use"`) ||
			strings.Contains(line, `"stop_reason"`) {
			t.Logf("Line %d: %s", i, line)
		}
	}
}

// TestToolCallsWithTextAfter tests tool calls followed by more text
func TestToolCallsWithTextAfter(t *testing.T) {
	testData, err := os.ReadFile("../../test/fixtures/morph_with_text_after_tools.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	input := bytes.NewReader(testData)
	var output bytes.Buffer

	err = TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Errorf("Transform failed: %v", err)
	}

	outputStr := output.String()

	// Check for tool use blocks
	toolUseCount := strings.Count(outputStr, `"type":"tool_use"`)
	t.Logf("Found %d tool_use blocks", toolUseCount)

	if toolUseCount == 0 {
		t.Error("CRITICAL: Expected tool_use blocks but found none - tools after text are being ignored!")
	}

	// Should NOT contain escaped XML in output
	if strings.Contains(outputStr, "\\u003c") {
		t.Error("CRITICAL: XML tags are being escaped and output as text instead of being parsed as tools!")
	}

	// Verify stop reason
	if !strings.Contains(outputStr, `"stop_reason":"tool_use"`) {
		t.Error("Stop reason should be tool_use when tools are present")
	}

	// Check for specific tools expected in this response
	if !strings.Contains(outputStr, `"name":"Read"`) {
		t.Error("Missing Read tool calls")
	}

	if !strings.Contains(outputStr, `"name":"Glob"`) {
		t.Error("Missing Glob tool calls")
	}

	if !strings.Contains(outputStr, `"name":"Write"`) {
		t.Error("Missing Write tool call")
	}

	// Print summary
	t.Logf("Total output length: %d bytes", len(outputStr))
	lines := strings.Split(outputStr, "\n")
	t.Logf("Total lines: %d", len(lines))
}

// TestTransformFromAbsolutePath tests transformation from an absolute file path
// and writes the output to client_response.txt in the project root
func TestTransformFromAbsolutePath(t *testing.T) {
	// Read the absolute path from environment variable or use default
	inputPath := os.Getenv("UPSTREAM_STREAM_FILE")
	if inputPath == "" {
		t.Skip("Set UPSTREAM_STREAM_FILE environment variable to test with custom file")
	}

	testData, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("Failed to read test data from %s: %v", inputPath, err)
	}

	input := bytes.NewReader(testData)
	var output bytes.Buffer

	err = TransformMorphToClaudeStream(input, "claude-sonnet-4-5", 0, &output, nil)
	if err != nil {
		t.Errorf("Transform failed: %v", err)
	}

	outputStr := output.String()

	// Write output to client_response.txt
	outputPath := "../../client_response.txt"
	if err := os.WriteFile(outputPath, []byte(outputStr), 0644); err != nil {
		t.Fatalf("Failed to write output to %s: %v", outputPath, err)
	}

	t.Logf("Successfully transformed %d bytes from %s", len(testData), inputPath)
	t.Logf("Output written to %s (%d bytes)", outputPath, len(outputStr))
}
