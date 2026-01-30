package parser

import "strings"

// BlockInfo contains information about a tool call block
type BlockInfo struct {
	StartIndex int
	TagType    string
}

// FindToolCallBlockAtEnd finds tool call block at the end of text
func FindToolCallBlockAtEnd(text string) *BlockInfo {
	trimmed := strings.TrimRight(text, " \t\n\r")

	// Find first occurrence of each tool call tag
	firstFunctionCalls := strings.Index(trimmed, "<function_calls>")
	firstTool := strings.Index(trimmed, "<tool>")
	firstTools := strings.Index(trimmed, "<tools>")

	// Find the earliest tag
	firstOpenTag := -1
	tagType := ""

	if firstFunctionCalls != -1 {
		firstOpenTag = firstFunctionCalls
		tagType = "function_calls"
	}
	if firstTool != -1 && (firstOpenTag == -1 || firstTool < firstOpenTag) {
		firstOpenTag = firstTool
		tagType = "tool"
	}
	if firstTools != -1 && (firstOpenTag == -1 || firstTools < firstOpenTag) {
		firstOpenTag = firstTools
		tagType = "tools"
	}

	if firstOpenTag == -1 {
		return nil
	}

	// Use stack to find matching close tag
	openTag := "<" + tagType + ">"
	closeTag := "</" + tagType + ">"

	depth := 0
	pos := firstOpenTag
	lastClosePos := -1

	for pos < len(trimmed) {
		nextOpen := strings.Index(trimmed[pos:], openTag)
		nextClose := strings.Index(trimmed[pos:], closeTag)

		if nextOpen == -1 && nextClose == -1 {
			break
		}

		if nextOpen != -1 && (nextClose == -1 || nextOpen < nextClose) {
			// Found open tag
			depth++
			pos = pos + nextOpen + len(openTag)
		} else {
			// Found close tag
			depth--
			if depth == 0 {
				lastClosePos = pos + nextClose + len(closeTag)
			}
			pos = pos + nextClose + len(closeTag)
		}
	}

	// Check if tool call is valid
	if lastClosePos != -1 {
		// Has complete open and close tags - this is valid
		// Accept tool calls even if there's text after them
	} else if depth > 0 {
		// Unclosed tool call, might be in progress (streaming)
		// Allow this case
	} else {
		// No valid tool call found
		return nil
	}

	return &BlockInfo{
		StartIndex: firstOpenTag,
		TagType:    tagType,
	}
}

// HasCompleteToolCall checks if text has complete tool call
func HasCompleteToolCall(text string) bool {
	trimmed := strings.TrimRight(text, " \t\n\r")
	return strings.HasSuffix(trimmed, "</function_calls>") ||
		strings.HasSuffix(trimmed, "</tool>") ||
		strings.HasSuffix(trimmed, "</tools>")
}

// HasIncompleteToolCall checks if text has incomplete tool call
func HasIncompleteToolCall(text string) bool {
	trimmed := strings.TrimRight(text, " \t\n\r")

	// Find first occurrence of each tool call tag
	firstFunctionCalls := strings.Index(trimmed, "<function_calls>")
	firstTool := strings.Index(trimmed, "<tool>")
	firstTools := strings.Index(trimmed, "<tools>")

	firstOpenTag := -1
	tagType := ""

	if firstFunctionCalls != -1 {
		firstOpenTag = firstFunctionCalls
		tagType = "function_calls"
	}
	if firstTool != -1 && (firstOpenTag == -1 || firstTool < firstOpenTag) {
		firstOpenTag = firstTool
		tagType = "tool"
	}
	if firstTools != -1 && (firstOpenTag == -1 || firstTools < firstOpenTag) {
		firstOpenTag = firstTools
		tagType = "tools"
	}

	if firstOpenTag == -1 {
		return false
	}

	// Check if there's a corresponding close tag
	closeTag := "</" + tagType + ">"
	closeIndex := strings.LastIndex(trimmed, closeTag)

	// If no close tag or close tag is before open tag, tool call is incomplete
	return closeIndex == -1 || closeIndex < firstOpenTag
}