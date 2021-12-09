// FROM https://github.com/ebuchman/go-shell-pipes

package sys

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/joomcode/errorx"
)

// Execute runs a command naked
func Execute(filename string, args ...string) (string, error) {
	return ExecuteEx(filename, "", false, args...)
}

var output string

// ExecuteEx runs a command and allows to set its working directory or pass the user environment to the command
func ExecuteEx(filename string, directoryName string, environment bool, args ...string) (string, error) {
	defer timeTrack(time.Now(), "execution")
	output = ""
	cmd := exec.Command(filename, args...)
	header := fmt.Sprintf("executing %s with %d args", filepath.Base(filename), len(args))
	if directoryName != "" {
		cmd.Dir = directoryName
		header += " in directory " + directoryName
	} else {
		header += " in current working directory (" + filepath.Dir(os.Args[0]) + ")"
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
		ExecuteOutputWriter(stdout)
	}()

	go func() {
		defer wg.Done()
		ExecuteOutputWriter(stderr)
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Errorf("command execution failed: %s", output)
		return output, errorx.Decorate(err, "failed to execute "+fullCommand)
	}
	// log.Tracef("combined output:%s%s", NewLine, output)
	return output, nil
}

func timeTrack(start time.Time, name string) {
	log.Debugf("%s took %ss", name, time.Since(start))
}

// ExecuteStreamingToStdout is used to silence command output
// in the ExecuteOutputWriter below
var ExecuteStreamingToStdout = false

// ExecuteOutputWriter is called the go routines above to stream to stdout
// TODO: parameterize output channel, etc
func ExecuteOutputWriter(r io.Reader) {
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
