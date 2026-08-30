module github.com/tx7do/kratos-bootstrap/ai/model

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/sashabaranov/go-openai v1.42.0
	github.com/tx7do/kratos-bootstrap/api v0.0.45
)

require google.golang.org/protobuf v1.36.12 // indirect
