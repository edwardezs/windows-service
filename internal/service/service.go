package service

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/nixpare/process"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"win-svc/internal/config"
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
	Cmd            *exec.Cmd
	Log            *lumberjack.Logger
}

func New(cfg config.WindowsService) *WindowsService {
	logPath := cfg.LogFilePath
	if logPath == "" || !filepath.IsAbs(logPath) {
		logPath = filepath.Join(filepath.Dir(cfg.ChildExecPath), logFileName)
	}

	return &WindowsService{
		Name:           cfg.Name,
		Description:    cfg.Description,
		ParentExecPath: cfg.ParentExecPath,
		ChildExecPath:  cfg.ChildExecPath,
		ChildExecArgs:  cfg.ChildExecArgs,
		Log: &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    cfg.LogFileMaxSizeMB,
			MaxBackups: cfg.LogFileMaxBackups,
			MaxAge:     cfg.LogFileMaxAgeDays,
		},
	}
}

func (w *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	defer w.Log.Close()

	processExited := make(chan error)
	if err := w.startProcess(processExited); err != nil {
		w.Log.Write([]byte(fmt.Sprintf("Failed to start process: %s\n", err.Error())))
		return
	}

	w.Log.Write([]byte("Process started\n"))
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				if err := w.stopProcess(processExited); err != nil {
					w.Log.Write([]byte(fmt.Sprintf("Failed to stop process: %s\n", err.Error())))
				}
				w.Log.Write([]byte("Process stopped\n"))
				break loop
			default:
				w.Log.Write([]byte(fmt.Sprintf("Unexpected control request #%d\n", c)))
			}
		case err := <-processExited:
			if err != nil {
				w.Log.Write([]byte(fmt.Sprintf("Process exited with error: %s, attempting restart\n", err.Error())))
			} else {
				w.Log.Write([]byte("Process exited, attempting restart\n"))
			}
			timeout := time.Now().Add(changeStateTimeout)
			for {
				if timeout.Before(time.Now()) {
					w.Log.Write([]byte("Timeout waiting for process to restart exceeded\n"))
					break loop
				}
				if err := w.startProcess(processExited); err != nil {
					w.Log.Write([]byte("Failed to start process, retrying\n"))
				} else {
					w.Log.Write([]byte("Process restarted\n"))
					break
				}
				time.Sleep(changeStateDelay)
			}
		}
	}

	changes <- svc.Status{State: svc.Stopped}
	return
}

func (w *WindowsService) Start() error {
	scm, err := mgr.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to service manager")
		return ErrFailedToConnectToServiceManager
	}
	defer scm.Disconnect()

	service, err := scm.OpenService(w.Name)
	if err != nil {
		log.Error().Err(err).Msgf("Service %s is not installed", w.Name)
		return ErrServiceNotExist
	}
	defer service.Close()

	if err := service.Start(); err != nil {
		log.Error().Err(err).Msg("Failed to start Windows service")
		return ErrFailedToStartService
	}
	log.Info().Msgf("Service %s started", w.Name)

	return nil
}

func (w *WindowsService) Run() {
	if err := svc.Run(w.Name, w); err != nil {
		w.Log.Write([]byte(fmt.Sprintf("Failed to start service: %s\n", err.Error())))
	}
}

func (w *WindowsService) Stop() error {
	scm, err := mgr.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to service manager")
		return ErrFailedToConnectToServiceManager
	}
	defer scm.Disconnect()

	service, err := scm.OpenService(w.Name)
	if err != nil {
		log.Error().Err(err).Msgf("Service %s is not installed", w.Name)
		return ErrServiceNotExist
	}
	defer service.Close()

	status, err := service.Control(svc.Stop)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send stop command to service %s", w.Name)
		return ErrFailedToSendStop
	}

	timeout := time.Now().Add(changeStateTimeout)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			log.Error().Msg("Timeout waiting for service to stop exceeded")
			return ErrStopTimeoutExceeded
		}
		time.Sleep(changeStateDelay)
		status, err = service.Query()
		if err != nil {
			log.Error().Err(err).Msg("Could not retrieve service status")
			return ErrFailedToGetServiceStatus
		}
	}
	log.Info().Msgf("Service %s stopped", w.Name)

	return nil
}

func (w *WindowsService) Install() error {
	scm, err := mgr.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to service manager")
		return ErrFailedToConnectToServiceManager
	}
	defer scm.Disconnect()

	service, err := scm.OpenService(w.Name)
	if err == nil {
		service.Close()
		log.Error().Msgf("Service %s already installed", w.Name)
		return ErrServiceAlreadyExist
	}

	service, err = scm.CreateService(w.Name, w.ParentExecPath, mgr.Config{
		DisplayName: w.Name,
		Description: w.Description,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create service %s", w.Name)
		return ErrFailedToCreateService
	}
	defer service.Close()

	log.Info().Msgf("Service %s installed", w.Name)

	return nil
}

func (w *WindowsService) Delete() error {
	scm, err := mgr.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to service manager")
		return ErrFailedToConnectToServiceManager
	}
	defer scm.Disconnect()

	service, err := scm.OpenService(w.Name)
	if err != nil {
		log.Error().Err(err).Msgf("Service %s is not installed", w.Name)
		return ErrServiceNotExist
	}
	defer service.Close()

	status, err := service.Query()
	if err != nil {
		log.Error().Err(err).Msg("Could not retrieve service status")
		return ErrFailedToGetServiceStatus
	}

	if status.State == svc.Running {
		log.Info().Msgf("Service %s is running, stopping", w.Name)
		timeout := time.Now().Add(changeStateTimeout)
		status, err := service.Control(svc.Stop)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to send stop command to service %s", w.Name)
			return ErrFailedToSendStop
		}
		for status.State != svc.Stopped {
			if timeout.Before(time.Now()) {
				log.Error().Msg("Timeout waiting for service to stop exceeded")
				return ErrStopTimeoutExceeded
			}
			time.Sleep(changeStateDelay)
			status, err = service.Query()
			if err != nil {
				log.Error().Err(err).Msg("Could not retrieve service status")
				return ErrFailedToGetServiceStatus
			}
		}
		log.Info().Msgf("Service %s stopped", w.Name)
	}

	if err = service.Delete(); err != nil {
		log.Error().Err(err).Msgf("Failed to delete service %s", w.Name)
		return ErrFailedToDeleteService
	}

	log.Info().Msgf("Service %s uninstalled", w.Name)

	return nil
}

func (w *WindowsService) startProcess(processExited chan error) error {
	w.Cmd = exec.Command(w.ChildExecPath, w.ChildExecArgs...)
	w.Cmd.Stdout = w.Log
	w.Cmd.Stderr = w.Log
	if err := w.Cmd.Start(); err != nil {
		return ErrFailedToStartService
	}

	go func() {
		processExited <- w.Cmd.Wait()
	}()

	return nil
}

func (w *WindowsService) stopProcess(processExited chan error) error {
	if err := process.StopProcess(w.Cmd.Process.Pid); err != nil {
		return ErrFailedToStopService
	}
	if err := <-processExited; err != nil {
		return err
	}

	return nil
}
