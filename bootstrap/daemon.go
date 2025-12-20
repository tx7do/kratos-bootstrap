package bootstrap

import (
	"fmt"
	"os"
)

// stripSlice 从字符串切片中移除指定元素
func stripSlice(slice []string, element string) []string {
	for i := 0; i < len(slice); {
		if slice[i] == element && i != len(slice)-1 {
			slice = append(slice[:i], slice[i+1:]...)
		} else if slice[i] == element && i == len(slice)-1 {
			slice = slice[:i]
		} else {
			i++
		}
	}
	return slice
}

// BeDaemon 将当前进程转为守护进程（尝试启动脱离的子进程并退出父进程）
func BeDaemon(arg string) {
	childArgs := stripSlice(os.Args, arg)
	cmd := subProcess(childArgs)
	if cmd == nil || cmd.Process == nil {
		// 启动失败，继续在当前进程运行
		return
	}
	fmt.Printf("[*] Daemon started in PID: %d\n", cmd.Process.Pid)
	os.Exit(0)
}
