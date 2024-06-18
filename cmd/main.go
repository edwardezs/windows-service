package main

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	win "golang.org/x/sys/windows/svc"

	"github.com/edwardezs/win-svc/pkg/cli"
	"github.com/edwardezs/win-svc/pkg/config"
	"github.com/edwardezs/win-svc/pkg/service"
)

const svcName = "Example Windows Service"

func main() {
	isWinSvc, err := win.IsWindowsService()
	if err != nil {
		log.Error().Err(err).Msg("Failed to determine if application is running as Windows service")
	}

	if isWinSvc {
		exePath, err := os.Executable()
		if err != nil {
			return
		}
		cfgPath := filepath.Join(filepath.Dir(exePath), cli.CfgFlag.Value)
		cfg, err := config.New(cfgPath)
		if err != nil {
			return
		}
		svc := service.New(cfg)
		svc.Run()
	}

	app := cli.New(svcName)
	if err := app.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("An error occurred while running the application")
	}
}
