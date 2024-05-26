package cli

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// ServiceCmd - cli-commands for running app as Windows service in background
// Usage:
//
//	./service.exe service install
//	./service.exe service start
//	./service.exe service stop
//	./service.exe service delete
//
// Note:  	admin rights are required to install/start/stop/remove app as Windows service
var ServiceCmd = cli.Command{
	Name:  "service",
	Usage: "send command to Windows service",
	Subcommands: []cli.Command{
		{
			Name:   "install",
			Usage:  "Install the service",
			Action: serviceInstallCmd,
		},
		{
			Name:   "start",
			Usage:  "Start the service",
			Action: serviceStartCmd,
		},
		{
			Name:   "stop",
			Usage:  "Stop the service",
			Action: serviceStopCmd,
		},
		{
			Name:   "delete",
			Usage:  "Delete the service",
			Action: serviceDeleteCmd,
		},
	},
}

func serviceStartCmd(ctx *cli.Context) error {
	if err := appCtx.svc.Start(); err != nil {
		return errors.Wrap(err, "service finished with error")
	}

	return nil
}

func serviceStopCmd(ctx *cli.Context) error {
	if err := appCtx.svc.Stop(); err != nil {
		return errors.Wrap(err, "service finished with error")
	}

	return nil
}

func serviceInstallCmd(ctx *cli.Context) error {
	if err := appCtx.svc.Install(); err != nil {
		return errors.Wrap(err, "failed to install service")
	}

	return nil
}

func serviceDeleteCmd(ctx *cli.Context) error {
	if err := appCtx.svc.Delete(); err != nil {
		return errors.Wrap(err, "failed to uninstall service")
	}

	return nil
}
