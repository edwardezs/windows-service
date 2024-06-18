package service

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/edwardezs/win-svc/pkg/config"
)

func New(cfg config.WindowsServiceConfig) *WindowsService {
	logPath := cfg.LogFilePath
	if logPath == "" || !filepath.IsAbs(logPath) {
		logPath = filepath.Join(filepath.Dir(cfg.ChildExecPath), logPath)
	}

	return &WindowsService{
		Name:           cfg.Name,
		Description:    cfg.Description,
		ParentExecPath: cfg.ParentExecPath,
		ChildExecPath:  cfg.ChildExecPath,
		ChildExecArgs:  cfg.ChildExecArgs,
		log: &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    cfg.LogFileMaxSizeMB,
			MaxBackups: cfg.LogFileMaxBackups,
			MaxAge:     cfg.LogFileMaxAgeDays,
		},
	}
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
		w.log.Write([]byte(fmt.Sprintf("Failed to start service: %s\n", err.Error())))
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
