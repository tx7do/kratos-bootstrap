package bootstrap

import (
	"github.com/spf13/cobra"
)

var (
	flags = NewCommandFlags()
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
	return &CommandFlags{
		Conf:       "../../configs",
		Env:        "dev",
		ConfigHost: "127.0.0.1:8500",
		ConfigType: "consul",
		Daemon:     false,
	}
}

// AddFlags 将 flags 绑定到传入的 cobra.Command（通常是 root command）。
func (f *CommandFlags) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&f.Conf, "conf", "c", f.Conf, "config path, eg: -conf ../../configs")
	cmd.PersistentFlags().StringVarP(&f.Env, "env", "e", f.Env, "runtime environment, eg: -env dev")
	cmd.PersistentFlags().StringVarP(&f.ConfigHost, "chost", "s", f.ConfigHost, "config server host, eg: -chost 127.0.0.1:8500")
	cmd.PersistentFlags().StringVarP(&f.ConfigType, "ctype", "t", f.ConfigType, "config server type, eg: -ctype consul")
	cmd.PersistentFlags().BoolVarP(&f.Daemon, "daemon", "d", f.Daemon, "run app as a daemon with -d or --daemon")
}

func (f *CommandFlags) Init() {
	if f.Daemon {
		BeDaemon("-d")
	}
}

// NewRootCmd 创建根命令并绑定命令行参数和执行函数。
func NewRootCmd(f *CommandFlags, runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "A microservice server application",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			f.Init()
		},
		RunE: runE,
	}
	f.AddFlags(cmd)
	return cmd
}
