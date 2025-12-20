//go:build !windows

package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Unix 平台的 subProcess 实现，使用 Setsid
func subProcess(args []string) *exec.Cmd {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err == nil {
		cmd.Stdin = devNull
		cmd.Stdout = devNull
		cmd.Stderr = devNull
	} else {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err = cmd.Start(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[-] Error starting daemon: %s\n", err)
		if devNull != nil {
			_ = devNull.Close()
		}
		return nil
	}
	if devNull != nil {
		_ = devNull.Close()
	}
	return cmd
}
