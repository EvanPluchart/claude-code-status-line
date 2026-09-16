//go:build !windows

package detach

import (
	"os/exec"
	"syscall"
)

// configure configures the command to run in its own session so it survives parent exit.
func configure(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
