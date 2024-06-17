package process

import (
	"io"
	"os"
)

func (p *Process) prepareStdin(stdin io.Reader) error {
	if stdin == nil {
		var err error
		p.in, err = p.Exec.StdinPipe()
		return err
	}

	if stdin == os.Stdin && p.Exec.SysProcAttr.HideWindow {
		inheritConsole(p.Exec.SysProcAttr, true)
	}
	
	p.Exec.Stdin = stdin
	return nil
}
