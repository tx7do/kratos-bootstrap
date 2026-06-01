package eino

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	jsonschema "github.com/eino-contrib/jsonschema"
)

// NewToolNode 创建一个工具执行节点，用于在 Graph/Chain 中执行工具调用。
// ToolsNode 接收包含 ToolCalls 的 AssistantMessage，执行对应的工具并返回结果消息。
//
// 用法示例：
//
//	toolsNode, _ := eino.NewToolNode(ctx, &compose.ToolsNodeConfig{
//	    Tools: []tool.BaseTool{searchTool, calcTool},
//	})
func NewToolNode(ctx context.Context, config *compose.ToolsNodeConfig) (*compose.ToolsNode, error) {
	return compose.NewToolNode(ctx, config)
}

// ToolsNodeOption 和 ToolsNodeConfig 透传。
type ToolsNodeConfig = compose.ToolsNodeConfig
type ToolsNodeOption = compose.ToolsNodeOption

// WithToolOption 设置工具调用的额外选项。
func WithToolOption(opts ...tool.Option) compose.ToolsNodeOption {
	return compose.WithToolOption(opts...)
}

// WithToolList 设置可调用的工具列表。
func WithToolList(tools ...tool.BaseTool) compose.ToolsNodeOption {
	return compose.WithToolList(tools...)
}

// --- ToolInfo 创建辅助 ---

// NewToolInfo 创建一个工具描述信息。
// name: 工具名称（需唯一）
// desc: 工具描述（告诉模型何时使用此工具）
// params: 工具参数定义（可选，无参工具传 nil）
func NewToolInfo(name, desc string, params *schema.ParamsOneOf) *schema.ToolInfo {
	return &schema.ToolInfo{
		Name:        name,
		Desc:        desc,
		ParamsOneOf: params,
	}
}

// NewParamsOneOfByParams 通过参数描述映射创建工具参数定义。
func NewParamsOneOfByParams(params map[string]*schema.ParameterInfo) *schema.ParamsOneOf {
	return schema.NewParamsOneOfByParams(params)
}

// NewParamsOneOfByJSONSchema 通过 JSON Schema 创建工具参数定义。
func NewParamsOneOfByJSONSchema(s *jsonschema.Schema) *schema.ParamsOneOf {
	return schema.NewParamsOneOfByJSONSchema(s)
}

// NewParameterInfo 创建一个参数描述。
func NewParameterInfo(typ schema.DataType, desc string, required bool, opts ...ParameterInfoOption) *schema.ParameterInfo {
	pi := &schema.ParameterInfo{
		Type:     typ,
		Desc:     desc,
		Required: required,
	}
	for _, opt := range opts {
		opt(pi)
	}
	return pi
}

// ParameterInfoOption 用于配置 ParameterInfo 的额外字段。
type ParameterInfoOption func(*schema.ParameterInfo)

// WithEnum 设置参数的枚举值（仅用于 string 类型）。
func WithEnum(enum []string) ParameterInfoOption {
	return func(pi *schema.ParameterInfo) {
		pi.Enum = enum
	}
}

// WithSubParams 设置子参数（仅用于 object 类型）。
func WithSubParams(subParams map[string]*schema.ParameterInfo) ParameterInfoOption {
	return func(pi *schema.ParameterInfo) {
		pi.SubParams = subParams
	}
}

// WithElemInfo 设置数组元素类型（仅用于 array 类型）。
func WithElemInfo(elem *schema.ParameterInfo) ParameterInfoOption {
	return func(pi *schema.ParameterInfo) {
		pi.ElemInfo = elem
	}
}
