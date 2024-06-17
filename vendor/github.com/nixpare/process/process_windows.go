package process

import (
	"fmt"
	"syscall"
)

// See https://learn.microsoft.com/en-us/windows/win32/procthread/process-creation-flags
const (
	create_new_console       = 0x00000010
	create_new_process_group = 0x00000200
)

func initSysProcAttr() syscall.SysProcAttr {
	var spa syscall.SysProcAttr
	inheritConsole(&spa, false)
	showWindow(&spa, false)
	return spa
}

// stop generates a CTRL+C signal
func (p *Process) stop() error {
	if !p.running {
		return fmt.Errorf("process %s not started", p.ExecName)
	}

	return StopProcess(p.Exec.Process.Pid)
}

func showWindow(spa *syscall.SysProcAttr, flag bool) {
	spa.HideWindow = !flag
}

func (p *Process) ShowWindow(flag bool) {
	showWindow(&p.SysProcAttr, flag)
}

func inheritConsole(spa *syscall.SysProcAttr, flag bool) {
	if flag {
		spa.CreationFlags = 0
	} else {
		spa.CreationFlags = create_new_console | create_new_process_group
	}
}

func (p *Process) InheritConsole(flag bool) {
	inheritConsole(&p.SysProcAttr, flag)
}
