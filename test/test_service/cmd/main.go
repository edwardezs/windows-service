package main

import (
	"os"
	"path/filepath"

	"win-svc/internal/config"
	"win-svc/internal/service"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		return
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
}
