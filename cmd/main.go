package main

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	win "golang.org/x/sys/windows/svc"

	"win-svc/internal/cli"
	"win-svc/internal/config"
	"win-svc/internal/service"
)

func main() {
	isWinSvc, err := win.IsWindowsService()
	if err != nil {
		log.Error().Err(err).Msg("Failed to determine if application is running as Windows service")
	}

	if isWinSvc {
		exePath, err := os.Executable()
		if err != nil {
			os.Exit(1)
		}
		cfgPath := filepath.Join(filepath.Dir(exePath), cli.CfgFlag.Value)
		cfg, err := config.New(cfgPath)
		if err != nil {
			os.Exit(1)
		}
		svc := service.New(cfg)
		svc.Run()
	}

	app := cli.New()
	if err := app.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("An error occurred while running the application")
	}
}