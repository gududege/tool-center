//go:build windows

package process

import (
	"os"
	"os/exec"
	"syscall"
)

func setCreateProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags |= syscall.CREATE_NEW_PROCESS_GROUP
}

func setHideWindow(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
	setCreateProcessGroup(cmd)
}

func sendInterrupt(p *os.Process) bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	attached, _, _ := kernel32.NewProc("AttachConsole").Call(uintptr(p.Pid))
	if attached == 0 {
		// 进程无控制台（GUI 环境），无法发信号
		return false
	}

	// CTRL_BREAK_EVENT 发到目标进程组（子进程及其后代），不影响 Wails 主进程
	// 注意：这里不用 CTRL_C_EVENT，因为 CREATE_NEW_PROCESS_GROUP 会禁用该组内
	// 所有进程的 CTRL_C_EVENT 处理，导致信号无法投递。CTRL_BREAK_EVENT 不受此限制。
	kernel32.NewProc("GenerateConsoleCtrlEvent").Call(
		uintptr(1),     // CTRL_BREAK_EVENT
		uintptr(p.Pid), // 子进程组 ID，仅投递到该组
	)
	kernel32.NewProc("FreeConsole").Call()
	return true
}
