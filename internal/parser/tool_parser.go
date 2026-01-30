package parser

import (
	"encoding/json"
	"opus-api/internal/types"
	"regexp"
	"sort"
	"strings"
)

type TagPosition struct {
	Type  string
	Index int
	Name  string
}

type ParseResult struct {
	ToolCalls     []types.ParsedToolCall
	RemainingText string
}

// 需要保持字符串类型的参数白名单（不进行 JSON 类型转换）
var stringOnlyParams = map[string]bool{
	"taskId": true, // TaskUpdate, TaskGet 等工具的任务 ID 必须是字符串
}

// shouldKeepAsString 检查参数是否应该保持为字符串类型
func shouldKeepAsString(paramName string) bool {
	return stringOnlyParams[paramName]
}

type NextToolCallResult struct {
	ToolCall    *types.ParsedToolCall
	EndPosition int
	Found       bool
}

func ParseNextToolCall(text string) NextToolCallResult {
	invokeStart := strings.Index(text, "<invoke name=\"")
	if invokeStart == -1 {
		return NextToolCallResult{Found: false}
	}

	depth := 0
	pos := invokeStart
	for pos < len(text) {
		nextOpen := strings.Index(text[pos:], "<invoke")
		nextClose := strings.Index(text[pos:], "</invoke>")

		if nextOpen == -1 && nextClose == -1 {
			break
		}

		if nextOpen != -1 && (nextClose == -1 || nextOpen < nextClose) {
			depth++
			pos = pos + nextOpen + 7
		} else {
			depth--
			if depth == 0 {
				endPos := pos + nextClose + 9
				invokeBlock := text[invokeStart:endPos]
				toolCalls := parseInvokeTags(invokeBlock)
				if len(toolCalls) > 0 {
					return NextToolCallResult{
						ToolCall:    &toolCalls[0],
						EndPosition: endPos,
						Found:       true,
					}
				}
				return NextToolCallResult{Found: false}
			}
			pos = pos + nextClose + 9
		}
	}

	return NextToolCallResult{Found: false}
}

func ParseToolCalls(text string) ParseResult {
	blockInfo := FindToolCallBlockAtEnd(text)
	if blockInfo == nil {
		return ParseResult{
			ToolCalls:     []types.ParsedToolCall{},
			RemainingText: text,
		}
	}
	remainingText := strings.TrimSpace(text[:blockInfo.StartIndex])
	toolCallBlock := text[blockInfo.StartIndex:]
	openTag := "<" + blockInfo.TagType + ">"
	closeTag := "</" + blockInfo.TagType + ">"
	openTagIndex := strings.Index(toolCallBlock, openTag)
	closeTagIndex := strings.LastIndex(toolCallBlock, closeTag)
	var innerContent string
	if closeTagIndex != -1 && closeTagIndex > openTagIndex {
		innerContent = toolCallBlock[openTagIndex+len(openTag) : closeTagIndex]
	} else {
		innerContent = toolCallBlock[openTagIndex+len(openTag):]
	}
	toolCalls := parseInvokeTags(innerContent)
	return ParseResult{
		ToolCalls:     toolCalls,
		RemainingText: remainingText,
	}
}

func parseInvokeTags(innerContent string) []types.ParsedToolCall {
	var toolCalls []types.ParsedToolCall
	invokeStartRegex := regexp.MustCompile(`<invoke name="([^"]+)">`)
	invokeEndRegex := regexp.MustCompile(`</invoke>`)
	var positions []TagPosition
	for _, match := range invokeStartRegex.FindAllStringSubmatchIndex(innerContent, -1) {
		positions = append(positions, TagPosition{
			Type:  "start",
			Index: match[0],
			Name:  innerContent[match[2]:match[3]],
		})
	}
	for _, match := range invokeEndRegex.FindAllStringIndex(innerContent, -1) {
		positions = append(positions, TagPosition{
			Type:  "end",
			Index: match[0],
		})
	}
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Index < positions[j].Index
	})
	depth := 0
	var currentInvoke *struct {
		Name       string
		StartIndex int
	}
	type InvokeBlock struct {
		Name  string
		Start int
		End   int
	}
	var topLevelInvokes []InvokeBlock
	for _, pos := range positions {
		if pos.Type == "start" {
			if depth == 0 {
				currentInvoke = &struct {
					Name       string
					StartIndex int
				}{Name: pos.Name, StartIndex: pos.Index}
			}
			depth++
		} else {
			depth--
			if depth == 0 && currentInvoke != nil {
				topLevelInvokes = append(topLevelInvokes, InvokeBlock{
					Name:  currentInvoke.Name,
					Start: currentInvoke.StartIndex,
					End:   pos.Index + 9,
				})
				currentInvoke = nil
			}
		}
	}
	if currentInvoke != nil && depth > 0 {
		topLevelInvokes = append(topLevelInvokes, InvokeBlock{
			Name:  currentInvoke.Name,
			Start: currentInvoke.StartIndex,
			End:   len(innerContent),
		})
	}
	for _, invoke := range topLevelInvokes {
		invokeContent := innerContent[invoke.Start:invoke.End]
		input := make(map[string]interface{})
		invokeTagEnd := strings.Index(invokeContent, ">") + 1
		paramsContent := invokeContent[invokeTagEnd:]
		paramStartRegex := regexp.MustCompile(`<parameter name="([^"]+)">`)
		paramEndRegex := regexp.MustCompile(`</parameter>`)
		var paramPositions []TagPosition
		for _, match := range paramStartRegex.FindAllStringSubmatchIndex(paramsContent, -1) {
			paramPositions = append(paramPositions, TagPosition{
				Type:  "start",
				Index: match[0],
				Name:  paramsContent[match[2]:match[3]],
			})
		}
		for _, match := range paramEndRegex.FindAllStringIndex(paramsContent, -1) {
			paramPositions = append(paramPositions, TagPosition{
				Type:  "end",
				Index: match[0],
			})
		}
		sort.Slice(paramPositions, func(i, j int) bool {
			return paramPositions[i].Index < paramPositions[j].Index
		})
		paramDepth := 0
		var currentParam *struct {
			Name       string
			StartIndex int
		}
		for _, pos := range paramPositions {
			if pos.Type == "start" {
				if paramDepth == 0 {
					tagEndIndex := strings.Index(paramsContent[pos.Index:], ">") + pos.Index + 1
					currentParam = &struct {
						Name       string
						StartIndex int
					}{Name: pos.Name, StartIndex: tagEndIndex}
				}
				paramDepth++
			} else {
				paramDepth--
				if paramDepth == 0 && currentParam != nil {
					value := paramsContent[currentParam.StartIndex:pos.Index]
					trimmedValue := strings.TrimSpace(value)

					// 检查是否在白名单中，白名单参数保持字符串类型
					if shouldKeepAsString(currentParam.Name) {
						input[currentParam.Name] = trimmedValue
					} else {
						// 尝试解析 JSON 格式的参数值（支持所有 JSON 类型）
						var parsed interface{}
						if err := json.Unmarshal([]byte(trimmedValue), &parsed); err == nil {
							// 解析成功，使用解析后的值（支持布尔值、数字、对象、数组等）
							input[currentParam.Name] = parsed
						} else {
							// 解析失败，保留原始字符串
							input[currentParam.Name] = trimmedValue
						}
					}
					currentParam = nil
				}
			}
		}
		if currentParam != nil && paramDepth > 0 {
			value := paramsContent[currentParam.StartIndex:]
			trimmedValue := strings.TrimSpace(value)

			// 检查是否在白名单中，白名单参数保持字符串类型
			if shouldKeepAsString(currentParam.Name) {
				input[currentParam.Name] = trimmedValue
			} else {
				// 尝试解析 JSON 格式的参数值（支持所有 JSON 类型）
				var parsed interface{}
				if err := json.Unmarshal([]byte(trimmedValue), &parsed); err == nil {
					// 解析成功，使用解析后的值（支持布尔值、数字、对象、数组等）
					input[currentParam.Name] = parsed
				} else {
					// 解析失败，保留原始字符串
					input[currentParam.Name] = trimmedValue
				}
			}
		}
		toolCalls = append(toolCalls, types.ParsedToolCall{Name: invoke.Name, Input: input})
	}
	return toolCalls
}
