package utils

import (
	"os"
	"os/exec"
	"runtime"
)

// GetShell 返回当前操作系统的 shell 和参数
func GetShell() (shell string, args []string) {
	if runtime.GOOS == "windows" {
		return "cmd", []string{}
	}

	// 优先使用环境变量中的 SHELL
	if envShell := os.Getenv("SHELL"); envShell != "" {
		return envShell, []string{}
	}

	// macOS 默认使用 zsh
	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("/bin/zsh"); err == nil {
			return "/bin/zsh", []string{}
		}
	}

	// Linux 默认使用 bash
	return "/bin/bash", []string{}
}

// GetShellCommand 返回执行命令的 shell 和参数
func GetShellCommand(command string) (shell string, args []string) {
	shell, _ = GetShell()
	if runtime.GOOS == "windows" {
		return shell, []string{"/c", command}
	}
	return shell, []string{"-c", command}
}

// NewShellCmd 创建一个交互式 shell 命令
func NewShellCmd() *exec.Cmd {
	shell, _ := GetShell()
	if runtime.GOOS == "windows" {
		return exec.Command(shell)
	}
	// Unix 系统使用 -i 启用交互模式
	return exec.Command(shell, "-i")
}

// NewShellCommandCmd 创建一个执行指定命令的 shell 命令
func NewShellCommandCmd(command string) *exec.Cmd {
	shell, args := GetShellCommand(command)
	return exec.Command(shell, args...)
}
