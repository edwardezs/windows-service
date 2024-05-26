package service

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/nixpare/process"
	"golang.org/x/sys/windows/svc"
)

const (
	logFileName = "service.log"

	changeStateTimeout = 10 * time.Second
	changeStateDelay   = 1 * time.Second
)

type WindowsService struct {
	Name           string
	Description    string
	ParentExecPath string
	ChildExecPath  string
	ChildExecArgs  []string
	cmd            *exec.Cmd
	log            *lumberjack.Logger
}

func (w *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	defer w.log.Close()

	processExited := make(chan error)
	if err := w.startProcess(processExited); err != nil {
		w.log.Write([]byte(fmt.Sprintf("Failed to start process: %s\n", err.Error())))
		return
	}

	w.log.Write([]byte("Process started\n"))
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				if err := w.stopProcess(processExited); err != nil {
					w.log.Write([]byte(fmt.Sprintf("Failed to stop process: %s\n", err.Error())))
				}
				w.log.Write([]byte("Process stopped\n"))
				break loop
			default:
				w.log.Write([]byte(fmt.Sprintf("Unexpected control request #%d\n", c)))
			}
		case err := <-processExited:
			if err != nil {
				w.log.Write([]byte(fmt.Sprintf("Process exited with error: %s, attempting restart\n", err.Error())))
			} else {
				w.log.Write([]byte("Process exited, attempting restart\n"))
			}
			timeout := time.Now().Add(changeStateTimeout)
			for {
				if timeout.Before(time.Now()) {
					w.log.Write([]byte("Timeout waiting for process to restart exceeded\n"))
					break loop
				}
				if err := w.startProcess(processExited); err != nil {
					w.log.Write([]byte("Failed to start process, retrying\n"))
				} else {
					w.log.Write([]byte("Process restarted\n"))
					break
				}
				time.Sleep(changeStateDelay)
			}
		}
	}

	changes <- svc.Status{State: svc.Stopped}
	return
}

func (w *WindowsService) startProcess(processExited chan error) error {
	w.cmd = exec.Command(w.ChildExecPath, w.ChildExecArgs...)
	w.cmd.Stdout = w.log
	w.cmd.Stderr = w.log
	if err := w.cmd.Start(); err != nil {
		return ErrFailedToStartService
	}

	go func() {
		processExited <- w.cmd.Wait()
	}()

	return nil
}

func (w *WindowsService) stopProcess(processExited chan error) error {
	if err := process.StopProcess(w.cmd.Process.Pid); err != nil {
		return ErrFailedToStopService
	}
	if err := <-processExited; err != nil {
		return err
	}

	return nil
}
