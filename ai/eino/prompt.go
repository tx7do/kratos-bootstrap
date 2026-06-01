package eino

import (
	einoPrompt "github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// FromMessages 根据消息模板列表创建 ChatTemplate。
// formatType 指定变量替换语法：
//   - schema.FString: Python 风格，如 "Hello, {name}!"
//   - schema.GoTemplate: Go 模板语法，如 "Hello, {{.name}}!"
//   - schema.Jinja2: Jinja2 语法，如 "Hello, {{ name }}!"
//
// 用法示例：
//
//	tpl := eino.FromMessages(schema.FString,
//	    &schema.Message{Role: schema.System, Content: "你是一个{role}助手。"},
//	    &schema.Message{Role: schema.User, Content: "{question}"},
//	)
//	messages, _ := tpl.Format(ctx, map[string]any{"role": "翻译", "question": "翻译这段话"})
func FromMessages(formatType schema.FormatType, templates ...schema.MessagesTemplate) *einoPrompt.DefaultChatTemplate {
	return einoPrompt.FromMessages(formatType, templates...)
}

// SystemMessage 创建一个系统角色的消息模板。
func SystemMessage(content string) *schema.Message {
	return &schema.Message{Role: schema.System, Content: content}
}

// UserMessage 创建一个用户角色的消息模板。
func UserMessage(content string) *schema.Message {
	return &schema.Message{Role: schema.User, Content: content}
}

// AssistantMessage 创建一个助手角色的消息模板。
func AssistantMessage(content string) *schema.Message {
	return &schema.Message{Role: schema.Assistant, Content: content}
}

// MessagesPlaceholder 创建一个消息占位符，用于在模板中插入动态消息列表。
// 例如在对话历史场景中，可以将历史消息注入到模板中。
//
// 用法示例：
//
//	tpl := eino.FromMessages(schema.FString,
//	    eino.SystemMessage("你是一个有用的助手。"),
//	    eino.MessagesPlaceholder("history", false),
//	    eino.UserMessage("{question}"),
//	)
func MessagesPlaceholder(variableName string, optional bool) schema.MessagesTemplate {
	return schema.MessagesPlaceholder(variableName, optional)
}
