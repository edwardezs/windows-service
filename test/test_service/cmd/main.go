package main

import (
	"os"
	"path/filepath"

	win "golang.org/x/sys/windows/svc"

	"win-svc/internal/config"
	"win-svc/internal/service"
)

func main() {
	isWinSvc, err := win.IsWindowsService()
	if err != nil {
		os.Exit(1)
	}
	if isWinSvc {
		exePath, err := os.Executable()
		if err != nil {
			os.Exit(1)
		}
		parentExecPath := filepath.Join(filepath.Dir(exePath), "test_service.exe")
		childExecPath := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(exePath))), "test_server/cmd/test_server.exe")
		service := service.New(config.WindowsService{
			Name:           "test_service",
			Description:    "Test Windows service",
			ParentExecPath: parentExecPath,
			ChildExecPath:  childExecPath,
			LogFilePath:    filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(exePath))), "test_server/cmd/test_service.log"),
		})
		service.Run()
		os.Exit(0)
	}
	os.Exit(0)
}
