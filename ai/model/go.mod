module github.com/tx7do/kratos-bootstrap/ai/model

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/sashabaranov/go-openai v1.41.2
	github.com/tx7do/kratos-bootstrap/api v0.0.41
)

require google.golang.org/protobuf v1.36.11 // indirect
