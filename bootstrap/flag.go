package bootstrap

import (
	"fmt"

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
	cmd.PersistentFlags().StringVarP(&f.ConfigHost, "chost", "ch", f.ConfigHost, "config server host, eg: -chost 127.0.0.1:8500")
	cmd.PersistentFlags().StringVarP(&f.ConfigType, "ctype", "ct", f.ConfigType, "config server type, eg: -ctype consul")
	cmd.PersistentFlags().BoolVarP(&f.Daemon, "daemon", "d", f.Daemon, "run app as a daemon with -d or --daemon")
}

func (f *CommandFlags) Init() {
	if f.Daemon {
		BeDaemon("-d")
	}

	ai := GetAppInfo()
	fmt.Printf("Application: %s\n", ai.Name)
	fmt.Printf("Version: %s\n", ai.Version)
	fmt.Printf("AppId: %s\n", ai.AppId)
	fmt.Printf("InstanceId: %s\n", ai.InstanceId)
	if len(ai.Metadata) > 0 {
		fmt.Println("Metadata:")
		for k, v := range ai.Metadata {
			fmt.Printf("  %s=%s\n", k, v)
		}
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
