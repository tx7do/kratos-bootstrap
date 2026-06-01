package langchaingo

import (
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

// NewOneShotAgent 创建一个 ReAct 风格的单次 Agent。
// Agent 会根据工具的名称和描述来决定是否调用工具，或直接给出最终答案。
func NewOneShotAgent(llm llms.Model, agentTools []tools.Tool, opts ...agents.Option) *agents.OneShotZeroAgent {
	return agents.NewOneShotAgent(llm, agentTools, opts...)
}

// NewConversationalAgent 创建一个对话式 Agent。
// 与 OneShotAgent 不同，它可以在使用工具的同时与用户进行自然对话。
func NewConversationalAgent(llm llms.Model, agentTools []tools.Tool, opts ...agents.Option) *agents.ConversationalAgent {
	return agents.NewConversationalAgent(llm, agentTools, opts...)
}

// NewOpenAIFunctionsAgent 创建一个基于 OpenAI Function Calling 的 Agent。
// 利用 OpenAI 的原生函数调用能力，相比 ReAct 解析更加稳定可靠。
func NewOpenAIFunctionsAgent(llm llms.Model, agentTools []tools.Tool, opts ...agents.Option) *agents.OpenAIFunctionsAgent {
	return agents.NewOpenAIFunctionsAgent(llm, agentTools, opts...)
}

// NewExecutor 创建一个 Agent 执行器，负责驱动 Agent 的推理-行动循环。
// 执行器会反复调用 Agent 的 Plan 方法，直到获得最终答案或达到最大迭代次数。
func NewExecutor(agent agents.Agent, opts ...agents.Option) *agents.Executor {
	return agents.NewExecutor(agent, opts...)
}

// NewOneShotExecutor 便捷方法：创建一个 OneShot Agent 并包装为 Executor。
func NewOneShotExecutor(llm llms.Model, agentTools []tools.Tool, opts ...agents.Option) *agents.Executor {
	agent := agents.NewOneShotAgent(llm, agentTools, opts...)
	return agents.NewExecutor(agent, opts...)
}

// NewConversationalExecutor 便捷方法：创建一个 Conversational Agent 并包装为 Executor。
func NewConversationalExecutor(llm llms.Model, agentTools []tools.Tool, opts ...agents.Option) *agents.Executor {
	agent := agents.NewConversationalAgent(llm, agentTools, opts...)
	return agents.NewExecutor(agent, opts...)
}

// NewOpenAIFunctionsExecutor 便捷方法：创建一个 OpenAI Functions Agent 并包装为 Executor。
func NewOpenAIFunctionsExecutor(llm llms.Model, agentTools []tools.Tool, opts ...agents.Option) *agents.Executor {
	agent := agents.NewOpenAIFunctionsAgent(llm, agentTools, opts...)
	return agents.NewExecutor(agent, opts...)
}
