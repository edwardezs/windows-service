//go:build !windows
package process

import (
	"os"
	"syscall"
)

func initSysProcAttr() syscall.SysProcAttr {
	return syscall.SysProcAttr{}
}

// stop sends a CTRL+C signal
func (p *Process) stop() error {
	return p.Exec.Process.Signal(os.Interrupt)
}
