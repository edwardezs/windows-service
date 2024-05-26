package service

import "github.com/pkg/errors"

var (
	// ErrServiceNotExist is hardcoded for proper error handling
	ErrServiceNotExist                 = errors.New("The specified service does not exist as an installed service.")
	ErrStopTimeoutExceeded             = errors.New("stop timeout exceeded")
	ErrServiceAlreadyExist             = errors.New("service already exists")
	ErrFailedToConnectToServiceManager = errors.New("failed to connect to service manager")
	ErrFailedToCreateService           = errors.New("failed to create service")
	ErrFailedToStartService            = errors.New("failed to start service")
	ErrFailedToStopService             = errors.New("failed to stop service")
	ErrFailedToDeleteService           = errors.New("failed to delete service")
	ErrFailedToRetrieveServiceStatus   = errors.New("failed to retrieve service status")
	ErrFailedToSendStop                = errors.New("failed to send stop command")
	ErrFailedToGetServiceStatus        = errors.New("failed to get service status")
)
