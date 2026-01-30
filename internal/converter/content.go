package converter

import (
	"encoding/json"
	"fmt"
	"opus-api/internal/types"
	"strings"
)

func ExtractTextFromContent(content interface{}) string {
	if content == nil {
		return ""
	}
	if str, ok := content.(string); ok {
		return str
	}
	if blocks, ok := content.([]types.ClaudeContentBlock); ok {
		var textParts []string
		for _, block := range blocks {
			switch b := block.(type) {
			case types.ClaudeContentBlockText:
				textParts = append(textParts, b.Text)
			case types.ClaudeContentBlockToolUse:
				var params []string
				for k, v := range b.Input {
					var valueStr string
					if str, ok := v.(string); ok {
						valueStr = str
					} else {
						jsonBytes, _ := json.Marshal(v)
						valueStr = string(jsonBytes)
					}
					paramTag := fmt.Sprintf("<parameter name=\"%s\">%s</parameter>", k, valueStr)
					params = append(params, paramTag)
				}
				xml := fmt.Sprintf("<function_calls>\n<invoke name=\"%s\">\n%s\n</invoke></function_calls>", b.Name, strings.Join(params, "\n"))
				textParts = append(textParts, xml)
			case types.ClaudeContentBlockToolResult:
				var resultContent string
				if str, ok := b.Content.(string); ok {
					resultContent = str
				} else {
					resultContent = ExtractTextFromContent(b.Content)
				}
				xml := fmt.Sprintf("<function_results>\n<result>\n<tool_use_id>%s</tool_use_id>\n<output>%s</output>\n</result>\n</function_results>", b.ToolUseID, resultContent)
				textParts = append(textParts, xml)
			}
		}
		return strings.Join(textParts, "\n")
	}
	return ""
}

func ExtractSystemText(system interface{}) string {
	if system == nil {
		return ""
	}
	if str, ok := system.(string); ok {
		return str
	}
	if messages, ok := system.([]types.ClaudeSystemMessage); ok {
		var texts []string
		for _, msg := range messages {
			texts = append(texts, msg.Text)
		}
		return strings.Join(texts, "\n")
	}
	return ""
}
