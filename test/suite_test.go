package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"win-svc/internal/config"
	"win-svc/internal/service"
)

const (
	parentExecPath = "test_service/cmd/test_service.exe"
	parentGoFile   = "test_service/cmd/main.go"
	childExecPath  = "test_server/cmd/test_server.exe"
	childGoFile    = "test_server/cmd/main.go"
	childURL       = "http://localhost:8080/hello"
	logFile        = "test_server/cmd/test_service.log"
)

var cfg = config.WindowsService{
	Name:        "test_service",
	Description: "Test Windows service",
}

type WindowsServiceTestSuite struct {
	suite.Suite
	svc *service.WindowsService
}

// `make test` for execution
// requeries built binaries of parent and child processes
// admin rights required as well
func TestSuiteWindowsService(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(WindowsServiceTestSuite))
}

func (s *WindowsServiceTestSuite) SetupSuite() {
	wd, err := os.Getwd()
	require.NoError(s.T(), err)
	cfg.ParentExecPath = filepath.Join(wd, parentExecPath)
	cfg.ChildExecPath = filepath.Join(wd, childExecPath)
	s.svc = service.New(cfg)
}

func (s *WindowsServiceTestSuite) TearDownSuite() {
	require.NoError(s.T(), os.Remove(parentExecPath))
	require.NoError(s.T(), os.Remove(childExecPath))
	require.NoError(s.T(), os.Remove(logFile))
}
