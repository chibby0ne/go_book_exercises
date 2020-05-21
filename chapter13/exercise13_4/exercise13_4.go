package exercise13_4

import (
	"io"
	"os/exec"
)

type writer struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

func NewWriter(out io.Writer) (io.WriteCloser, error) {
	cmd := exec.Cmd{
		Path:   "/bin/bzip2",
		Args:   []string{"-9"},
		Stdout: out,
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	// start the command but doesn't wait for it to finish, so it will start
	// and hang until some input (stdin) is given
	cmd.Start()
	w := &writer{cmd: &cmd, stdin: stdin}
	return w, nil
}

func (w *writer) Write(data []byte) (int, error) {
	return w.stdin.Write(data)
}

func (w *writer) Close() error {
	if err := w.stdin.Close(); err != nil {
		return err
	}
	if err := w.cmd.Wait(); err != nil {
		return err
	}
	return nil
}
