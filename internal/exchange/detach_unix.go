//go:build !windows

package exchange

import (
	"os/exec"
	"syscall"
)

// detach configures the command to run in its own session so it survives parent exit.
func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
