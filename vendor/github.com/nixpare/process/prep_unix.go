//go:build !windows
package process

import (
	"io"
)

func (p *Process) prepareStdin(stdin io.Reader) error {
	if stdin == nil {
		var err error
		p.in, err = p.Exec.StdinPipe()
		return err
	}
	
	p.Exec.Stdin = stdin
	return nil
}
