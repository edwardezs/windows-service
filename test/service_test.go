package test

import (
	"net/http"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/stretchr/testify/require"

	"win-svc/internal/service"
)

const (
	installDelay = 500 * time.Millisecond
	startDelay   = 500 * time.Millisecond
	restartDelay = 500 * time.Millisecond
	stopDelay    = 500 * time.Millisecond
	deleteDelay  = 500 * time.Millisecond
)

func (s *WindowsServiceTestSuite) TestExecution() {
	require.NoError(s.T(), s.svc.Install())

	err := s.svc.Install()
	require.Error(s.T(), err)
	require.Equal(s.T(), err, service.ErrServiceAlreadyExist)

	time.Sleep(installDelay)

	require.NoError(s.T(), s.svc.Start())

	time.Sleep(startDelay)

	resp, err := http.Get(childURL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	require.NoError(s.T(), s.svc.Stop())

	time.Sleep(stopDelay)

	_, err = http.Get(childURL)
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "connectex: No connection could be made")

	require.NoError(s.T(), s.svc.Delete())

	time.Sleep(deleteDelay)

	err = s.svc.Delete()
	require.Error(s.T(), err)
	require.Equal(s.T(), err, service.ErrServiceNotExist)
}

func (s *WindowsServiceTestSuite) TestChildProcessKill() {
	require.NoError(s.T(), s.svc.Install())

	time.Sleep(installDelay)

	require.NoError(s.T(), s.svc.Start())

	time.Sleep(startDelay)

	killCmd := exec.Command("cmd", "/C", "TASKKILL", "/F", "/IM", filepath.Base(childExecPath))
	require.NoError(s.T(), killCmd.Run())

	time.Sleep(restartDelay)

	resp, err := http.Get(childURL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	require.NoError(s.T(), s.svc.Stop())

	time.Sleep(stopDelay)

	require.NoError(s.T(), s.svc.Delete())
}
