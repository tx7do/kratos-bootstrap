package langchaingo

import (
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

// NewLLMChain 创建一个 LLM 链，将提示模板和 LLM 组合在一起。
// 这是最基础的链，用于将用户的输入通过 Prompt 模板格式化后传递给 LLM 生成回答。
func NewLLMChain(llm llms.Model, prompt prompts.FormatPrompter, opts ...chains.ChainCallOption) *chains.LLMChain {
	return chains.NewLLMChain(llm, prompt, opts...)
}

// NewConversationChain 创建一个对话链，内置对话记忆管理。
// 自动维护对话历史，适合多轮对话场景。
func NewConversationChain(llm llms.Model, mem schema.Memory) chains.LLMChain {
	if mem == nil {
		mem = memory.NewConversationBuffer()
	}
	return chains.NewConversation(llm, mem)
}

// NewSequentialChain 创建一个顺序链，将多个链按顺序串联执行。
// 前一个链的输出作为后一个链的输入。
func NewSequentialChain(c []chains.Chain, inputKeys []string, outputKeys []string, opts ...chains.SequentialChainOption) (*chains.SequentialChain, error) {
	return chains.NewSequentialChain(c, inputKeys, outputKeys, opts...)
}

// NewSimpleSequentialChain 创建一个简单顺序链，所有链必须只有单一输入和单一输出。
func NewSimpleSequentialChain(c []chains.Chain) (*chains.SimpleSequentialChain, error) {
	return chains.NewSimpleSequentialChain(c)
}

// NewStuffDocumentsChain 创建一个文档填充链，将所有文档填充到单个提示中。
func NewStuffDocumentsChain(llmChain *chains.LLMChain) chains.StuffDocuments {
	return chains.NewStuffDocuments(llmChain)
}

// LoadStuffSummarization 加载一个摘要链，将所有文档内容填充到提示中进行摘要。
func LoadStuffSummarization(llm llms.Model) chains.StuffDocuments {
	return chains.LoadStuffSummarization(llm)
}

// LoadRefineSummarization 加载一个精炼摘要链，通过迭代精炼方式生成摘要。
func LoadRefineSummarization(llm llms.Model) chains.RefineDocuments {
	return chains.LoadRefineSummarization(llm)
}

// LoadMapReduceSummarization 加载一个 MapReduce 摘要链，先分块摘要再合并。
func LoadMapReduceSummarization(llm llms.Model) chains.MapReduceDocuments {
	return chains.LoadMapReduceSummarization(llm)
}
