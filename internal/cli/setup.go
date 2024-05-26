package cli

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"win-svc/internal/config"
	"win-svc/internal/service"
)

// CfgFlag is the cli-flag used for parsing configuration of the Windows service
// Usage: 		./service.exe <-config> <FLAG_VALUE> ...
// Required:	false
// Default: 	service.config.json
var CfgFlag = &cli.StringFlag{
	Name:  "config",
	Value: "service.config.json", // default
	Usage: "Configuration file",
}

var appCtx AppContext

type AppContext struct {
	svc *service.WindowsService
	cfg config.WindowsService
}

// SetupService sets up the AppContext for the Windows service
func SetupService(ctx *cli.Context) (err error) {
	appCtx.cfg, err = config.New(ctx.String(CfgFlag.Name))
	if err != nil {
		return errors.Wrap(err, "failed to load config for Windows service")
	}
	appCtx.svc = service.New(appCtx.cfg)

	return nil
}
