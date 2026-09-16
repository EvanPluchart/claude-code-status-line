// Package detach spawns background helper processes that outlive the statusline render.
//
// The render path must stay within its execution budget, so any refresh that needs
// the network runs in a detached child process writing to an on-disk cache.
package detach

import (
	"fmt"
	"os"
	"os/exec"
)

// Spawn starts the current executable with args in a detached process and releases it.
func Spawn(args ...string) error {
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate executable: %w", err)
	}

	cmd := exec.Command(self, args...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	configure(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", args[0], err)
	}

	// Release the child so it outlives this process.
	return cmd.Process.Release()
}
