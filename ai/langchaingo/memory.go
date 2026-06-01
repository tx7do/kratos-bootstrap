package langchaingo

import (
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

// NewChatMessageHistory 创建一个内存中的聊天消息历史记录。
// 用于手动管理对话消息的存储和检索。
func NewChatMessageHistory(opts ...memory.ChatMessageHistoryOption) *memory.ChatMessageHistory {
	return memory.NewChatMessageHistory(opts...)
}

// NewConversationBuffer 创建一个完整对话缓冲记忆。
// 保存所有对话历史，不做任何裁剪。适合对话轮次较少的场景。
func NewConversationBuffer(opts ...memory.ConversationBufferOption) *memory.ConversationBuffer {
	return memory.NewConversationBuffer(opts...)
}

// NewConversationWindowBuffer 创建一个窗口对话缓冲记忆。
// 只保留最近 N 轮对话（默认 5 轮），自动丢弃更早的历史。
// 适合需要控制上下文长度的场景。
func NewConversationWindowBuffer(windowSize int, opts ...memory.ConversationBufferOption) *memory.ConversationWindowBuffer {
	return memory.NewConversationWindowBuffer(windowSize, opts...)
}

// NewConversationTokenBuffer 创建一个基于 Token 数量的对话缓冲记忆。
// 当对话历史超过指定的 Token 限制时，自动从最老的开始裁剪。
// 需要传入 LLM 用于计算 Token 数量。
func NewConversationTokenBuffer(llm llms.Model, maxTokenLimit int, opts ...memory.ConversationBufferOption) *memory.ConversationTokenBuffer {
	return memory.NewConversationTokenBuffer(llm, maxTokenLimit, opts...)
}

// NewSimpleMemory 创建一个空记忆（不存储任何内容）。
// 用于不需要记忆的 Chain 或 Agent。
func NewSimpleMemory() memory.Simple {
	return memory.NewSimple()
}

// --- 常用 Memory Option 透传 ---

// WithChatHistory 设置自定义的聊天历史存储后端。
func WithChatHistory(chatHistory schema.ChatMessageHistory) memory.ConversationBufferOption {
	return memory.WithChatHistory(chatHistory)
}

// WithReturnMessages 设置是否以消息对象形式返回（而非字符串）。
func WithReturnMessages(returnMessages bool) memory.ConversationBufferOption {
	return memory.WithReturnMessages(returnMessages)
}

// WithMemoryKey 设置记忆变量的键名（默认为 "history"）。
func WithMemoryKey(memoryKey string) memory.ConversationBufferOption {
	return memory.WithMemoryKey(memoryKey)
}

// WithInputKey 设置输入键名。
func WithInputKey(inputKey string) memory.ConversationBufferOption {
	return memory.WithInputKey(inputKey)
}

// WithOutputKey 设置输出键名。
func WithOutputKey(outputKey string) memory.ConversationBufferOption {
	return memory.WithOutputKey(outputKey)
}

// WithHumanPrefix 设置用户消息前缀（默认为 "Human"）。
func WithHumanPrefix(humanPrefix string) memory.ConversationBufferOption {
	return memory.WithHumanPrefix(humanPrefix)
}

// WithAIPrefix 设置 AI 消息前缀（默认为 "AI"）。
func WithAIPrefix(aiPrefix string) memory.ConversationBufferOption {
	return memory.WithAIPrefix(aiPrefix)
}
