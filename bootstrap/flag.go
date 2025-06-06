package bootstrap

import (
	"flag"

	"github.com/tx7do/kratos-bootstrap/utils"
)

// CommandFlags 命令传参
type CommandFlags struct {
	Conf       string // 引导配置文件路径，默认为：../../configs
	Env        string // 开发环境：dev、debug……
	ConfigHost string // 远程配置服务端地址
	ConfigType string // 远程配置服务端类型
	Daemon     bool   // 是否转为守护进程
}

func NewCommandFlags() *CommandFlags {
	f := &CommandFlags{
		Conf:       "",
		Env:        "",
		ConfigHost: "",
		ConfigType: "",
		Daemon:     false,
	}

	f.defineFlag()

	return f
}

func (f *CommandFlags) defineFlag() {
	flag.StringVar(&f.Conf, "conf", "../../configs", "config path, eg: -conf ../../configs")
	flag.StringVar(&f.Env, "env", "dev", "runtime environment, eg: -env dev")
	flag.StringVar(&f.ConfigHost, "chost", "127.0.0.1:8500", "config server host, eg: -chost 127.0.0.1:8500")
	flag.StringVar(&f.ConfigType, "ctype", "consul", "config server host, eg: -ctype consul")
	flag.BoolVar(&f.Daemon, "d", false, "run app as a daemon with -d=true.")
}

func (f *CommandFlags) Init() {
	flag.Parse()

	if f.Daemon {
		utils.BeDaemon("-d")
	}
}
