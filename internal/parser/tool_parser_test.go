package parser

import (
	"testing"
)

func TestTaskIdAsString(t *testing.T) {
	// 测试 taskId 参数应该保持为字符串类型，即使值看起来像数字
	text := `<function_calls>
<invoke name="TaskUpdate">
<parameter name="taskId">1</parameter>
<parameter name="status">completed</parameter>
</invoke>
</function_calls>`

	result := ParseToolCalls(text)

	if len(result.ToolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(result.ToolCalls))
	}

	toolCall := result.ToolCalls[0]
	if toolCall.Name != "TaskUpdate" {
		t.Errorf("Expected tool name 'TaskUpdate', got '%s'", toolCall.Name)
	}

	// 检查 taskId 是否为字符串类型
	taskId, ok := toolCall.Input["taskId"]
	if !ok {
		t.Fatal("taskId parameter not found")
	}

	taskIdStr, ok := taskId.(string)
	if !ok {
		t.Errorf("taskId should be string type, got %T with value %v", taskId, taskId)
	}

	if taskIdStr != "1" {
		t.Errorf("Expected taskId to be '1', got '%s'", taskIdStr)
	}

	// 检查 status 参数应该保持为字符串（因为它本来就是字符串）
	status, ok := toolCall.Input["status"]
	if !ok {
		t.Fatal("status parameter not found")
	}

	statusStr, ok := status.(string)
	if !ok {
		t.Errorf("status should be string type, got %T", status)
	}

	if statusStr != "completed" {
		t.Errorf("Expected status to be 'completed', got '%s'", statusStr)
	}
}

func TestOtherNumericParams(t *testing.T) {
	// 测试其他数字参数应该被解析为数字类型
	text := `<function_calls>
<invoke name="SomeTool">
<parameter name="count">42</parameter>
<parameter name="enabled">true</parameter>
</invoke>
</function_calls>`

	result := ParseToolCalls(text)

	if len(result.ToolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(result.ToolCalls))
	}

	toolCall := result.ToolCalls[0]

	// count 应该是数字类型
	count, ok := toolCall.Input["count"]
	if !ok {
		t.Fatal("count parameter not found")
	}

	countFloat, ok := count.(float64)
	if !ok {
		t.Errorf("count should be float64 type, got %T", count)
	}

	if countFloat != 42 {
		t.Errorf("Expected count to be 42, got %f", countFloat)
	}

	// enabled 应该是布尔类型
	enabled, ok := toolCall.Input["enabled"]
	if !ok {
		t.Fatal("enabled parameter not found")
	}

	enabledBool, ok := enabled.(bool)
	if !ok {
		t.Errorf("enabled should be bool type, got %T", enabled)
	}

	if !enabledBool {
		t.Errorf("Expected enabled to be true, got %v", enabledBool)
	}
}