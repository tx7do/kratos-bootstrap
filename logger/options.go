package logger

type Type string

const (
	Std     Type = "std"
	Fluent  Type = "fluent"
	Logrus  Type = "logrus"
	Zap     Type = "zap"
	Aliyun  Type = "aliyun"
	Tencent Type = "tencent"
)
