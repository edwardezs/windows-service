package test

import (
	"net/http"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/stretchr/testify/require"

	"win-svc/internal/service"
)

func (s *WindowsServiceTestSuite) TestExecution() {
	require.NoError(s.T(), s.svc.Install())

	err := s.svc.Install()
	require.Error(s.T(), err)
	require.Equal(s.T(), err, service.ErrServiceAlreadyExist)

	require.NoError(s.T(), s.svc.Start())

	time.Sleep(1 * time.Second)

	resp, err := http.Get(childURL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	require.NoError(s.T(), s.svc.Stop())

	_, err = http.Get(childURL)
	require.Error(s.T(), err)

	require.NoError(s.T(), s.svc.Delete())

	err = s.svc.Delete()
	require.Error(s.T(), err)
	require.Equal(s.T(), err, service.ErrServiceNotExist)
}

func (s *WindowsServiceTestSuite) TestChildProcessKill() {
	require.NoError(s.T(), s.svc.Install())

	require.NoError(s.T(), s.svc.Start())

	time.Sleep(1 * time.Second)

	killCmd := exec.Command("cmd", "/C", "TASKKILL", "/F", "/IM", filepath.Base(childExecPath))
	require.NoError(s.T(), killCmd.Run())

	time.Sleep(5 * time.Second) // Waiting for the child process to restart

	resp, err := http.Get(childURL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	require.NoError(s.T(), s.svc.Stop())

	require.NoError(s.T(), s.svc.Delete())
}
