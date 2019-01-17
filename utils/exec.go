// FROM https://github.com/ebuchman/go-shell-pipes

package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debugf("%s took %ss", name, elapsed)
}

func Execute(filename string, args ...string) (string, error) {
	return ExecuteEx(filename, "", false, args...)
}

func ExecuteEx(filename string, directoryName string, environment bool, args ...string) (string, error) {
	defer timeTrack(time.Now(), "execution")
	fullCommand := fmt.Sprintf("%s %s", filename, strings.Join(args[:], " "))
	log.Debugf("executing with %d args%s%s%sin directory %s with environment: %t", len(args), NewLine, fullCommand, NewLine, directoryName, environment)
	cmd := exec.Command(filename, args...)
	if directoryName != "" {
		cmd.Dir = directoryName
	}
	if environment {
		cmd.Env = os.Environ()
	}
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("command execution failed: %s", string(bytes))
		return string(bytes), errorx.Decorate(err, "failed to execute "+fullCommand)
	}
	out := string(bytes)
	log.Tracef("combined output:%s%s", NewLine, out)
	return out, nil
}

func execute(cmd *exec.Cmd, timeout <-chan time.Time) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = mergeWriter(cmd.Stdout, &stdout)
	cmd.Stderr = mergeWriter(cmd.Stderr, &stderr)

	log.Infof("Execute '%s %s'", cmd.Path, strings.Join(cmd.Args[1:], " ")) // skip arg[0] as it is printed separately
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting %v:\nCommand stdout:\n%v\nstderr:\n%v\nerror:\n%v", cmd, stdout.String(), stderr.String(), err)
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()
	log.Debugf("Entering select")
	select {
	case err := <-errCh:
		log.Debugf("Received on errCh: %v", err)
		if err != nil {
			return stderr.String(), err
		}
	case <-timeout:
		if err := cmd.Process.Kill(); err != nil {
			return "", err
		}
		return "", fmt.Errorf(
			"timed out waiting for command %v:\nCommand stdout:\n%v\nstderr:\n%v",
			cmd.Args, stdout.String(), stderr.String())
	}
	log.Debugf("stderr: %s", stderr.String())
	return stdout.String(), nil
}

func mergeWriter(other io.Writer, buf io.Writer) io.Writer {
	if other != nil {
		return io.MultiWriter(other, buf)
	}
	return buf
}

// Convert a shell command with a series of pipes into
// correspondingly piped list of *exec.Cmd
// If an arg has spaces, this will fail
func RunString(s string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	sp := strings.Split(s, "|")
	cmds := make([]*exec.Cmd, len(sp))
	// create the commands
	for i, c := range sp {
		cs := strings.Split(strings.TrimSpace(c), " ")
		cmd := cmdFromStrings(cs)
		cmds[i] = cmd
	}

	cmds = AssemblePipes(cmds, nil, buf)
	if err := RunCmds(cmds); err != nil {
		return "", err
	}

	b := buf.Bytes()
	return string(b), nil
}

func cmdFromStrings(cs []string) *exec.Cmd {
	if len(cs) == 1 {
		return exec.Command(cs[0])
	} else if len(cs) == 2 {
		return exec.Command(cs[0], cs[1])
	}
	return exec.Command(cs[0], cs[1:]...)
}

// Convert sequence of tokens into commands,
// using "|" as a delimiter
func RunStrings(tokens ...string) (string, error) {
	if len(tokens) == 0 {
		return "", nil
	}
	buf := bytes.NewBuffer([]byte{})
	cmds := []*exec.Cmd{}
	args := []string{}
	// accumulate tokens until a |
	for _, t := range tokens {
		if t != "|" {
			args = append(args, t)
		} else {
			cmds = append(cmds, cmdFromStrings(args))
			args = []string{}
		}
	}
	cmds = append(cmds, cmdFromStrings(args))
	cmds = AssemblePipes(cmds, nil, buf)
	if err := RunCmds(cmds); err != nil {
		return "", fmt.Errorf("%s; %s", err.Error(), string(buf.Bytes()))
	}

	b := buf.Bytes()
	return string(b), nil
}

// Pipe stdout of each command into stdin of next
func AssemblePipes(cmds []*exec.Cmd, stdin io.Reader, stdout io.Writer) []*exec.Cmd {
	cmds[0].Stdin = stdin
	cmds[0].Stderr = stdout
	// assemble pipes
	for i, c := range cmds {
		if i < len(cmds)-1 {
			cmds[i+1].Stdin, _ = c.StdoutPipe()
			cmds[i+1].Stderr = stdout
		} else {
			c.Stdout = stdout
			c.Stderr = stdout
		}
	}
	return cmds
}

// Run series of piped commands
func RunCmds(cmds []*exec.Cmd) error {
	// start processes in descending order
	for i := len(cmds) - 1; i > 0; i-- {
		if err := cmds[i].Start(); err != nil {
			return err
		}
	}
	// run the first process
	if err := cmds[0].Run(); err != nil {
		return err
	}
	// wait on processes in ascending order
	for i := 1; i < len(cmds); i++ {
		if err := cmds[i].Wait(); err != nil {
			return err
		}
	}
	return nil
}
