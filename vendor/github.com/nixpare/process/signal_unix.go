//go:build !windows
package process

import "os"

// StopProcess sends an os.Interrupt to the process
func StopProcess(PID int) error {
	p, err := os.FindProcess(PID)
	if err != nil {
		return err
	}

	return p.Signal(os.Interrupt)
}