package eino

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// --- Lambda 封装 ---

// Lambda 类型透传。
type Lambda = compose.Lambda

// LambdaOption 透传。
type LambdaOption = compose.LambdaOpt

// InvokableLambda 创建一个仅支持 Invoke 模式的 Lambda 节点。
// 用于在 Chain/Graph 中插入自定义同步逻辑。
//
// 用法示例：
//
//	lambda := eino.InvokableLambda(func(ctx context.Context, input string) (string, error) {
//	    return strings.ToUpper(input), nil
//	})
//	chain.AppendLambda(lambda)
func InvokableLambda[I, O any](fn func(ctx context.Context, input I) (O, error), opts ...compose.LambdaOpt) *compose.Lambda {
	return compose.InvokableLambda(fn, opts...)
}

// StreamableLambda 创建一个支持 Stream 模式的 Lambda 节点。
func StreamableLambda[I, O any](fn func(ctx context.Context, input I) (*schema.StreamReader[O], error), opts ...compose.LambdaOpt) *compose.Lambda {
	return compose.StreamableLambda(fn, opts...)
}

// CollectableLambda 创建一个支持 Collect 模式的 Lambda 节点。
func CollectableLambda[I, O any](fn func(ctx context.Context, input *schema.StreamReader[I]) (O, error), opts ...compose.LambdaOpt) *compose.Lambda {
	return compose.CollectableLambda(fn, opts...)
}

// TransformableLambda 创建一个支持 Transform 模式的 Lambda 节点。
func TransformableLambda[I, O any](fn func(ctx context.Context, input *schema.StreamReader[I]) (*schema.StreamReader[O], error), opts ...compose.LambdaOpt) *compose.Lambda {
	return compose.TransformableLambda(fn, opts...)
}

// AnyLambda 创建一个支持多种模式的 Lambda 节点。
// 传入 nil 表示不支持对应模式，但至少需要一个非 nil。
func AnyLambda[I, O, TOption any](
	invoke compose.Invoke[I, O, TOption],
	stream compose.Stream[I, O, TOption],
	collect compose.Collect[I, O, TOption],
	transform compose.Transform[I, O, TOption],
	opts ...compose.LambdaOpt,
) (*compose.Lambda, error) {
	return compose.AnyLambda(invoke, stream, collect, transform, opts...)
}

// ToList 创建一个将单个输入转换为切片的 Lambda 节点。
// 常用于将 ChatModel 输出的单条消息转换为消息列表。
//
// 用法示例：
//
//	toListLambda := eino.ToList[*schema.Message]()
//	chain.AppendLambda(toListLambda)
func ToList[I any](opts ...compose.LambdaOpt) *compose.Lambda {
	return compose.ToList[I](opts...)
}

// --- Parallel 封装 ---

// Parallel 类型透传。
type Parallel = compose.Parallel

// NewParallel 创建一个并行结构，用于在 Chain 中同时执行多个节点。
// 并行节点的输出类型为 map[string]any，其中 key 由 Add* 方法的 outputKey 参数指定。
//
// 用法示例：
//
//	parallel := eino.NewParallel()
//	parallel.AddLambda("upper", upperLambda)
//	parallel.AddLambda("lower", lowerLambda)
//	chain.AppendParallel(parallel)
func NewParallel() *compose.Parallel {
	return compose.NewParallel()
}

// --- ChainBranch 封装 ---

// ChainBranch 类型透传。
type ChainBranch = compose.ChainBranch

// GraphBranchCondition 类型透传。
type GraphBranchCondition[T any] = compose.GraphBranchCondition[T]

// NewChainBranch 创建一个条件分支，根据条件函数的返回值选择执行哪个节点。
//
// 用法示例：
//
//	branch := eino.NewChainBranch(func(ctx context.Context, input string) (string, error) {
//	    if len(input) > 100 {
//	        return "summarize", nil
//	    }
//	    return "direct", nil
//	})
//	branch.AddLambda("summarize", summarizeLambda)
//	branch.AddLambda("direct", directLambda)
//	chain.AppendBranch(branch)
func NewChainBranch[T any](cond compose.GraphBranchCondition[T]) *compose.ChainBranch {
	return compose.NewChainBranch(cond)
}

// --- Workflow 封装 ---

// NewWorkflow 创建一个工作流，支持依赖管理和复杂的数据流转。
// Workflow 比 Graph 更高级，支持通过 AddInput/AddDependency 声明节点间的数据依赖关系，
// 而不是手动添加 Edge。
//
// 用法示例：
//
//	wf := eino.NewWorkflow[map[string]any, *schema.Message]()
//	promptNode := wf.AddChatTemplateNode("prompt", chatTemplate)
//	modelNode := wf.AddChatModelNode("model", chatModel)
//	modelNode.AddInput("prompt")
//	wf.AddEnd("model")
//	r, _ := wf.Compile(ctx)
func NewWorkflow[I, O any](opts ...compose.NewGraphOption) *compose.Workflow[I, O] {
	return compose.NewWorkflow[I, O](opts...)
}

// WorkflowNode 类型透传。
type WorkflowNode = compose.WorkflowNode

// --- Runnable 封装 ---

// Runnable 类型透传，表示编译后的可执行 Chain/Graph/Workflow。
type Runnable[I, O any] = compose.Runnable[I, O]

// --- StateGraph 辅助 ---

// WithGenLocalState 注册一个函数来生成每次运行时的本地状态。
// 状态可以在节点间共享，适用于需要跨节点传递中间数据的场景。
//
// 用法示例：
//
//	type MyState struct {
//	    History []*schema.Message
//	}
//	graph := eino.NewGraph[map[string]any, *schema.Message](
//	    eino.WithGenLocalState(func(ctx context.Context) *MyState {
//	        return &MyState{}
//	    }),
//	)
func WithGenLocalState[S any](gls func(ctx context.Context) S) compose.NewGraphOption {
	return compose.WithGenLocalState(gls)
}

// --- Graph Node Option 透传 ---

// GraphAddNodeOpt 透传。
type GraphAddNodeOpt = compose.GraphAddNodeOpt

// GraphCompileOption 透传。
type GraphCompileOption = compose.GraphCompileOption
