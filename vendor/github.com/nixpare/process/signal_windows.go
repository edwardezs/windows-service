package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const ctrl_c_command = "--nix-send-ctrl"

func init() {
	if len(os.Args) < 3 || os.Args[1] != ctrl_c_command {
		return
	}

	PID := -1
	_, err := fmt.Sscanf(os.Args[2], "%d", &PID)
	if err != nil || PID <= 0 {
		fmt.Print(err)
		os.Exit(1)
	}

	if PID <= 0 {
		fmt.Printf("invalid PID: %d", PID)
		os.Exit(1)
	}

	err = stopProcessThread(PID)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	os.Exit(0)
}

// StopProcess simulates a CTRL+C signal to the Process
func StopProcess(PID int) error {
	child := exec.Command(os.Args[0], ctrl_c_command, fmt.Sprint(PID))
	
	b, err := child.CombinedOutput()
	if err != nil {
		return fmt.Errorf("CTRL-C thread error: %w", err)
	}
	if len(b) > 0 {
		return fmt.Errorf("CTRL-C thread error: %s", string(b))
	}

	return nil
}

func stopProcessThread(PID int) error {
	err := FreeConsole()
	if err != nil {
		return err
	}

	err = AttachConsole(PID)
	if err != nil {
		return err
	}

	err = RemoveConsoleCtrlHandler()
	if err != nil {
		return err
	}

	err = GenerateConsoleCtrlCEvent()
	if err != nil {
		return err
	}

	return nil
}

// LoadKernel32 loads the kernel32 dll that can be used to find all sort of windows APIs
func LoadKernel32() (*syscall.DLL, error) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return nil, fmt.Errorf("loadDLL: %w", e)
	}

	return d, nil
}

// GenerateConsoleCtrlCEvent generates a CTRL+C event that spreads to all the processes attached to
// the current underlying console
func GenerateConsoleCtrlCEvent() error {
	d, e := LoadKernel32()
	if e != nil {
		return e
	}

	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return fmt.Errorf("findProc: %w", e)
	}

	r, _, e := p.Call(syscall.CTRL_C_EVENT, 0)
	if r == 0 {
		return fmt.Errorf("generateConsoleCtrlEvent: %w", e)
	}

	return nil
}

// GenerateConsoleCtrlBreakEvent generates a CTRL+C event that spreads to all the processes attached to
// the current underlying console
func GenerateConsoleCtrlBreakEvent() error {
	d, e := LoadKernel32()
	if e != nil {
		return e
	}

	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return fmt.Errorf("findProc: %w", e)
	}

	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, 0)
	if r == 0 {
		return fmt.Errorf("generateConsoleCtrlEvent: %w", e)
	}

	return nil
}


func setConsoleCtrlHandler(flag bool) error {
	d, e := LoadKernel32()
	if e != nil {
		return e
	}

	p, e := d.FindProc("SetConsoleCtrlHandler")
	if e != nil {
		return fmt.Errorf("findProc: %w", e)
	}

	a := 0
	if flag {
		a = 1
	}

	r, _, e := p.Call(0, uintptr(a))
	if r == 0 {
		return fmt.Errorf("setConsoleCtrlHandler: %w", e)
	}

	return nil
}

// RemoveConsoleCtrlHandler makes the calling process not a target of the CTRL+C event
func RemoveConsoleCtrlHandler() error {
	return setConsoleCtrlHandler(true)
}

// RestoreConsoleCtrlHandler reverts the effect of RemoveConsoleCtrlHandler function and makes
// the calling process a target of the CTRL+C event
func RestoreConsoleCtrlHandler() error {
	return setConsoleCtrlHandler(false)
}

// FreeConsole detatches the calling process from the underlying console
func FreeConsole() error {
	d, e := LoadKernel32()
	if e != nil {
		return e
	}

	p, e := d.FindProc("FreeConsole")
	if e != nil {
		return fmt.Errorf("findProc: %v", e)
	}

	r, _, e := p.Call(0)
	if r == 0 {
		return fmt.Errorf("freeConsole: %w", e)
	}

	return nil
}

// AllocConsole creates a new console for the calling process (it must not already have one)
func AllocConsole() error {
	d, e := LoadKernel32()
	if e != nil {
		return e
	}

	p, e := d.FindProc("AllocConsole")
	if e != nil {
		return fmt.Errorf("findProc: %v", e)
	}

	r, _, e := p.Call(0)
	if r == 0 {
		return fmt.Errorf("allocConsole: %w", e)
	}

	return nil
}

// AttachConsole attaches to the same console as the process with the given PID (it must not already have one)
func AttachConsole(pid int) error {
	d, e := LoadKernel32()
	if e != nil {
		return e
	}

	p, e := d.FindProc("AttachConsole")
	if e != nil {
		return fmt.Errorf("findProc: %v", e)
	}

	r, _, e := p.Call(uintptr(pid))
	if r == 0 {
		return fmt.Errorf("attachConsole: %w", e)
	}

	return nil
}

// Makes che calling process immune to CTRL+C events and generates one of them.
// You may call RestoreConsoleCtrlHandler to revert the change after all child
// processes have exited
func SendCtrlC() error {
	if e := RemoveConsoleCtrlHandler(); e != nil {
		return e
	}

	if e := GenerateConsoleCtrlCEvent(); e != nil {
		return e
	}

	return nil
}
