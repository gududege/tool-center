//go:build !windows

package process

import (
	"os"
	"os/exec"
)

func setCreateProcessGroup(cmd *exec.Cmd) {}

func setHideWindow(cmd *exec.Cmd) {}

func sendInterrupt(p *os.Process) bool {
	p.Signal(os.Interrupt)
	return true
}
