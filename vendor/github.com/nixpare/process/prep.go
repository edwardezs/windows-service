package process

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/nixpare/broadcaster"
)

func (p *Process) prepareStdout(stdout io.Writer) error {
	p.stdOutErrWG.Add(1)
	outPipe, err := p.Exec.StdoutPipe()
	if err != nil {
		return nil
	}

	go func() {
		defer func() {
			p.outBc.Close()
			p.stdOutErrWG.Done()
		}()
		pipeOutput(p.outBc, outPipe, stdout, "stdout")
	}()
	return nil
}

func (p *Process) prepareStderr(stderr io.Writer) error {
	p.stdOutErrWG.Add(1)
	errPipe, err := p.Exec.StderrPipe()
	if err != nil {
		return nil
	}

	go func() {
		defer func() {
			p.errBc.Close()
			p.stdOutErrWG.Done()
		}()
		pipeOutput(p.errBc, errPipe, stderr, "stderr")
	}()
	return nil
}

func pipeOutput(bc *broadcaster.BufBroadcaster[[]byte], r io.ReadCloser, w io.Writer, pipeID string) {
	pipeR, pipeW := io.Pipe()
	var buf [1024]byte

	br := bufio.NewReader(pipeR)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			line, err := br.ReadBytes('\n')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					line = []byte(fmt.Sprintf("broken %s pipe: %v", pipeID, err))
				}
			}

			if len(line) > 0 {
				if line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
				}

				bc.Send(line)
			}

			if err != nil {
				break
			}
		}
	}()

	for {
		n, err := r.Read(buf[:])
		b := buf[:n]
		if err != nil {
			if !errors.Is(err, io.EOF) {
				b = append(b, []byte(fmt.Sprintf("broken %s pipe: %v", pipeID, err))...)
			}
		}

		if len(b) > 0 {
			pipeW.Write(b)
			if w != nil && w != dev_null {
				w.Write(b)
			}
		}

		if err != nil {
			pipeW.CloseWithError(err)
			break
		}
	}

	wg.Wait()
}
