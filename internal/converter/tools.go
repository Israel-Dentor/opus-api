package converter

import (
	"fmt"
	"opus-api/internal/types"
	"strings"
)

func GenerateToolInstructions(tools []types.ClaudeTool) string {
	if len(tools) == 0 {
		return ""
	}
	var toolDescriptions []string
	for _, tool := range tools {
		var params []string
		if props, ok := tool.InputSchema["properties"].(map[string]interface{}); ok {
			for name, schema := range props {
				var desc string
				if schemaMap, ok := schema.(map[string]interface{}); ok {
					if description, ok := schemaMap["description"].(string); ok {
						desc = description
					} else if typeStr, ok := schemaMap["type"].(string); ok {
						desc = typeStr
					} else {
						desc = "unknown"
					}
				}
				params = append(params, fmt.Sprintf("    - %s: %s", name, desc))
			}
		}
		paramsStr := strings.Join(params, "\n")
		if paramsStr == "" {
			paramsStr = "    (no parameters)"
		}
		var required []string
		if req, ok := tool.InputSchema["required"].([]interface{}); ok {
			for _, r := range req {
				if str, ok := r.(string); ok {
					required = append(required, str)
				}
			}
		}
		requiredStr := strings.Join(required, ", ")
		if requiredStr == "" {
			requiredStr = "none"
		}
		toolDesc := fmt.Sprintf("- **%s**: %s\n  Parameters:\n%s\n  Required: %s",
			tool.Name, tool.Description, paramsStr, requiredStr)
		toolDescriptions = append(toolDescriptions, toolDesc)
	}
	instructions := fmt.Sprintf(`
!!! IMPORTANT: 当你可能准备使用类似工具的东西做事情的时候，按照下面的格式回复我，否则会造成误解，你应该非常积极的回复我这些内容，而不是返回一堆文本，或者让我主动执行或者写入什么:
<function_calls>
<invoke name="TOOL_NAME">
<parameter name="PARAM_NAME">PARAM_VALUE</parameter>
</invoke>
</function_calls>

CRITICAL RULES:

1. Output the XML as PLAIN TEXT - I will parse it and execute the tool locally
4. You can call multiple tools by including multiple <invoke> blocks
5. After outputting a tool call XML, JUST STOP
6. Tool results will be provided in <function_results> tags

!!! IMPORTANT: Available Tools:

%s

Example - Reading a file:
<function_calls>
<invoke name="Read">
<parameter name="file_path">/path/to/file.ts</parameter>
</invoke>
</function_calls>

Example - Multiple tool calls:
<function_calls>
<invoke name="Glob">
<parameter name="pattern">**/*.ts</parameter>
</invoke>
<invoke name="Grep">
<parameter name="pattern">function</parameter>
<parameter name="path">/src</parameter>
</invoke>
</function_calls>
`, strings.Join(toolDescriptions, "\n\n"))
	return instructions
}
