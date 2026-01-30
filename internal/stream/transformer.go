package stream

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"opus-api/internal/parser"
	"opus-api/internal/tokenizer"
	"opus-api/internal/types"
	"strings"
)

// TransformMorphToClaudeStream transforms MorphLLM SSE stream to Claude SSE stream
func TransformMorphToClaudeStream(morphStream io.Reader, model string, inputTokens int, writer io.Writer, onChunk func(string)) error {
	scanner := bufio.NewScanner(morphStream)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024) // Increase buffer size

	messageID := "msg_" + generateUUID()
	hasStarted := false
	contentBlockStarted := false
	contentBlockClosed := false
	messageDeltaSent := false
	toolCallsEmitted := false
	fullText := ""
	contentBlockIndex := 0
	buffer := NewTextBuffer()
	nativeToolCalls := []types.ParsedToolCall{}

	emitSSE := func(event string, data interface{}) {
		sseData := FormatSSE(event, data)
		if onChunk != nil {
			onChunk(sseData)
		}
		writer.Write([]byte(sseData))
	}

	emitToolCall := func(toolCall types.ParsedToolCall) {
		// Close current text block if open
		if contentBlockStarted && !contentBlockClosed {
			emitSSE("content_block_stop", ContentBlockStopEvent{
				Type:  "content_block_stop",
				Index: contentBlockIndex,
			})
			contentBlockClosed = true
		}
		contentBlockIndex++

		toolUseID := "toolu_" + generateShortUUID()

		emitSSE("content_block_start", ContentBlockStartEvent{
			Type:  "content_block_start",
			Index: contentBlockIndex,
			ContentBlock: ToolUseContentBlock{
				Type:  "tool_use",
				ID:    toolUseID,
				Name:  toolCall.Name,
				Input: toolCall.Input,
			},
		})

		emitSSE("content_block_delta", ContentBlockDeltaEvent{
			Type:  "content_block_delta",
			Index: contentBlockIndex,
			Delta: InputJSONDelta{
				Type:        "input_json_delta",
				PartialJSON: mustMarshalJSON(toolCall.Input),
			},
		})

		emitSSE("content_block_stop", ContentBlockStopEvent{
			Type:  "content_block_stop",
			Index: contentBlockIndex,
		})

		contentBlockIndex++
		toolCallsEmitted = true
	}

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		dataStr := strings.TrimPrefix(line, "data: ")
		dataStr = strings.TrimSpace(dataStr)

		if dataStr == "[DONE]" {
			// Calculate output tokens from accumulated text
			outputTokens := tokenizer.CountTokens(fullText)

			// Handle [DONE]
			if toolCallsEmitted {
				if !messageDeltaSent {
					emitSSE("message_delta", MessageDeltaEvent{
						Type: "message_delta",
						Delta: map[string]interface{}{
							"stop_reason":   "tool_use",
							"stop_sequence": nil,
						},
						Usage: map[string]int{"output_tokens": outputTokens},
					})
					messageDeltaSent = true
				}
				emitSSE("message_stop", MessageStopEvent{Type: "message_stop"})
				continue
			}

			// Check for native tool calls (backup)
			if len(nativeToolCalls) > 0 {
				for _, toolCall := range nativeToolCalls {
					emitToolCall(toolCall)
				}
				if !messageDeltaSent {
					emitSSE("message_delta", MessageDeltaEvent{
						Type: "message_delta",
						Delta: map[string]interface{}{
							"stop_reason":   "tool_use",
							"stop_sequence": nil,
						},
						Usage: map[string]int{"output_tokens": outputTokens},
					})
					messageDeltaSent = true
				}
				emitSSE("message_stop", MessageStopEvent{Type: "message_stop"})
				continue
			}

			// No tool calls, flush remaining text
			if !buffer.IsEmpty() {
				buffer.FlushAll(func(text string) {
					emitSSE("content_block_delta", ContentBlockDeltaEvent{
						Type:  "content_block_delta",
						Index: contentBlockIndex,
						Delta: TextDelta{Type: "text_delta", Text: text},
					})
				})
			}

			// Close text content block if open
			if contentBlockStarted && !contentBlockClosed {
				emitSSE("content_block_stop", ContentBlockStopEvent{
					Type:  "content_block_stop",
					Index: contentBlockIndex,
				})
				contentBlockClosed = true
			}

			// Send message_delta if not sent
			if !messageDeltaSent {
				emitSSE("message_delta", MessageDeltaEvent{
					Type: "message_delta",
					Delta: map[string]interface{}{
						"stop_reason":   "end_turn",
						"stop_sequence": nil,
					},
					Usage: map[string]int{"output_tokens": outputTokens},
				})
				messageDeltaSent = true
			}

			emitSSE("message_stop", MessageStopEvent{Type: "message_stop"})
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
			continue
		}

		dataType, _ := data["type"].(string)

		switch dataType {
		case "start":
			if !hasStarted {
				hasStarted = true
				emitSSE("message_start", MessageStartEvent{
					Type: "message_start",
					Message: MessageStart{
						ID:           messageID,
						Type:         "message",
						Role:         "assistant",
						Content:      []interface{}{},
						Model:        model,
						StopReason:   nil,
						StopSequence: nil,
						Usage:        map[string]int{"input_tokens": inputTokens, "output_tokens": 0},
					},
				})
			}

		case "text-start":
			if !contentBlockStarted {
				contentBlockStarted = true
				contentBlockClosed = false
				emitSSE("content_block_start", ContentBlockStartEvent{
					Type:         "content_block_start",
					Index:        contentBlockIndex,
					ContentBlock: TextContentBlock{Type: "text", Text: ""},
				})
			} else if contentBlockClosed {
				contentBlockIndex++
				contentBlockClosed = false
				emitSSE("content_block_start", ContentBlockStartEvent{
					Type:         "content_block_start",
					Index:        contentBlockIndex,
					ContentBlock: TextContentBlock{Type: "text", Text: ""},
				})
			}

		case "text-delta":
			delta, _ := data["delta"].(string)
			fullText += delta

			// If tool calls already emitted, ignore subsequent text
			if toolCallsEmitted {
				continue
			}

			// If content block closed, reopen it
			if contentBlockClosed {
				contentBlockIndex++
				contentBlockClosed = false
				emitSSE("content_block_start", ContentBlockStartEvent{
					Type:         "content_block_start",
					Index:        contentBlockIndex,
					ContentBlock: TextContentBlock{Type: "text", Text: ""},
				})
			}

			buffer.Add(delta)

			// Stream processing: emit tool calls one by one as they complete
			for {
				result := parser.ParseNextToolCall(fullText)
				if !result.Found {
					break
				}

				// Output text before tool call
				textBefore := fullText[:strings.Index(fullText, "<invoke")]
				if textBefore != "" && !buffer.ToolCallDetected {
					emitSSE("content_block_delta", ContentBlockDeltaEvent{
						Type:  "content_block_delta",
						Index: contentBlockIndex,
						Delta: TextDelta{Type: "text_delta", Text: textBefore},
					})
				}
				buffer.Clear()
				buffer.ToolCallDetected = true

				// Emit single tool call immediately
				emitToolCall(*result.ToolCall)

				// Remove processed part from fullText
				fullText = fullText[result.EndPosition:]
			}

			// Check for incomplete tool call
			if parser.HasIncompleteToolCall(fullText) {
				buffer.ToolCallDetected = true
				buffer.Clear()
			} else if !buffer.ToolCallDetected {
				// No tool call, output text normally
				buffer.FlushSafeText(func(text string) {
					emitSSE("content_block_delta", ContentBlockDeltaEvent{
						Type:  "content_block_delta",
						Index: contentBlockIndex,
						Delta: TextDelta{Type: "text_delta", Text: text},
					})
				})
			}

		case "text-end":
			// text-end indicates current text segment ended
			result := parser.ParseToolCalls(fullText)

			if len(result.ToolCalls) == 0 {
				// No tool calls, output all remaining text
				if !buffer.IsEmpty() {
					buffer.FlushAll(func(text string) {
						emitSSE("content_block_delta", ContentBlockDeltaEvent{
							Type:  "content_block_delta",
							Index: contentBlockIndex,
							Delta: TextDelta{Type: "text_delta", Text: text},
						})
					})
				}
			}

		case "finish-step":
			// MorphLLM finish-step indicates a step completed
			// Check for tool calls at step boundary
			result := parser.ParseToolCalls(fullText)
			if len(result.ToolCalls) > 0 && !toolCallsEmitted {
				// Output remaining text before tool calls
				if result.RemainingText != "" && !buffer.ToolCallDetected {
					emitSSE("content_block_delta", ContentBlockDeltaEvent{
						Type:  "content_block_delta",
						Index: contentBlockIndex,
						Delta: TextDelta{Type: "text_delta", Text: result.RemainingText},
					})
				}
				buffer.Clear()
				buffer.ToolCallDetected = true

				// Emit tool calls
				for _, toolCall := range result.ToolCalls {
					emitToolCall(toolCall)
				}
			} else if !buffer.IsEmpty() && !buffer.ToolCallDetected {
				// No tool calls, flush remaining text
				buffer.FlushAll(func(text string) {
					emitSSE("content_block_delta", ContentBlockDeltaEvent{
						Type:  "content_block_delta",
						Index: contentBlockIndex,
						Delta: TextDelta{Type: "text_delta", Text: text},
					})
				})
			}

		case "start-step":
			// MorphLLM start-step indicates new step started
			// No special handling needed

		case "finish":
			result := parser.ParseToolCalls(fullText)

			finishReason, _ := data["finishReason"].(string)
			if len(result.ToolCalls) == 0 && finishReason != "tool-calls" && !messageDeltaSent {
				stopReason := "end_turn"
				if finishReason != "" && finishReason != "stop" {
					stopReason = finishReason
				}
				outputTokens := tokenizer.CountTokens(fullText)
				emitSSE("message_delta", MessageDeltaEvent{
					Type: "message_delta",
					Delta: map[string]interface{}{
						"stop_reason":   stopReason,
						"stop_sequence": nil,
					},
					Usage: map[string]int{"output_tokens": outputTokens},
				})
				messageDeltaSent = true
			}

		case "tool-input-error":
			// Capture MorphLLM native tool calls (when tool unavailable)
			toolName, _ := data["toolName"].(string)
			input, _ := data["input"].(map[string]interface{})
			if toolName != "" && input != nil {
				nativeToolCalls = append(nativeToolCalls, types.ParsedToolCall{
					Name:  toolName,
					Input: input,
				})
				buffer.ToolCallDetected = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
}

func generateShortUUID() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func mustMarshalJSON(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}