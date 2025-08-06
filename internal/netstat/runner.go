// internal/netstat/runner.go
package netstat

import (
	"bufio"
	"fmt"
	"os/exec"
)

// Runner executes netstat command and returns output
type Runner struct {
	args []string
}

// NewRunner creates a new Runner instance
func NewRunner(args []string) *Runner {
	return &Runner{args: args}
}

// Run executes the netstat command and returns a scanner for reading output
// The cleanup function should be called when done reading
func (r *Runner) Run() (*bufio.Scanner, func(), error) {
	binaryPath, err := exec.LookPath("netstat")
	if err != nil {
		return nil, nil, fmt.Errorf("netstat not found: %w", err)
	}

	cmd := exec.Command(binaryPath, r.args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start netstat: %w", err)
	}

	cleanup := func() {
		cmd.Wait()
	}

	scanner := bufio.NewScanner(stdout)
	return scanner, cleanup, nil
}
