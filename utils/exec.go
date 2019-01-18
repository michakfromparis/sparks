// FROM https://github.com/ebuchman/go-shell-pipes

package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
)

var ExecuteStreamingToStdout = false

func Execute(filename string, args ...string) (string, error) {
	return ExecuteEx(filename, "", false, args...)
}

var output string

func ExecuteEx(filename string, directoryName string, environment bool, args ...string) (string, error) {
	defer timeTrack(time.Now(), "execution")
	output = ""
	cmd := exec.Command(filename, args...)
	header := fmt.Sprintf("executing %s with %d args", filepath.Base(filename), len(args))
	if directoryName != "" {
		cmd.Dir = directoryName
		header += " in directory " + directoryName
	} else {
		header += " in default directory (" + filepath.Dir(os.Args[0]) + ")"
	}
	if environment {
		cmd.Env = os.Environ()
		header += " and user environment"
	}
	fullCommand := fmt.Sprintf("%s %s", filename, strings.Join(args[:], " "))
	log.Debugf("%s%s%s", header, NewLine, fullCommand)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		writer(stdout)
	}()

	go func() {
		defer wg.Done()
		writer(stderr)
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Errorf("command execution failed: %s", output)
		return output, errorx.Decorate(err, "failed to execute "+fullCommand)
	}
	log.Tracef("combined output:%s%s", NewLine, output)
	return output, nil
}

func timeTrack(start time.Time, name string) {
	log.Debugf("%s took %ss", name, time.Since(start))
}

func writer(r io.Reader) {
	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output += string(buf[0:n])
			if ExecuteStreamingToStdout {
				fmt.Print(string(buf[0:n]))
			}
		}
		if err != nil {
			break
		}
	}
}

// Convert a shell command with a series of pipes into
// correspondingly piped list of *exec.Cmd
// If an arg has spaces, this will fail
func ExecutePipe(s string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	sp := strings.Split(s, "|")
	cmds := make([]*exec.Cmd, len(sp))
	// create the commands
	for i, c := range sp {
		cs := strings.Split(strings.TrimSpace(c), " ")
		cmd := cmdFromString(cs)
		cmds[i] = cmd
	}

	cmds = AssemblePipes(cmds, nil, buf)
	if err := runCmds(cmds); err != nil {
		return "", err
	}

	b := buf.Bytes()
	return string(b), nil
}

func cmdFromString(cs []string) *exec.Cmd {
	if len(cs) == 1 {
		return exec.Command(cs[0])
	} else if len(cs) == 2 {
		return exec.Command(cs[0], cs[1])
	}
	return exec.Command(cs[0], cs[1:]...)
}

// Convert sequence of tokens into commands,
// using "|" as a delimiter
func ExecutePipes(tokens ...string) (string, error) {
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
			cmds = append(cmds, cmdFromString(args))
			args = []string{}
		}
	}
	cmds = append(cmds, cmdFromString(args))
	cmds = AssemblePipes(cmds, nil, buf)
	if err := runCmds(cmds); err != nil {
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
func runCmds(cmds []*exec.Cmd) error {
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
