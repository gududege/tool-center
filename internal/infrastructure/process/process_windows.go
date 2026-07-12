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
	setCreateProcessGroup(cmd)
	cmd.SysProcAttr.HideWindow = true
}

func sendInterrupt(p *os.Process) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	attached, _, _ := kernel32.NewProc("AttachConsole").Call(uintptr(p.Pid))
	kernel32.NewProc("GenerateConsoleCtrlEvent").Call(
		uintptr(1),       // CTRL_BREAK_EVENT
		uintptr(p.Pid),   // process group = child's PID
	)
	if attached != 0 {
		kernel32.NewProc("FreeConsole").Call()
	}
}
