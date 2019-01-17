// FROM https://github.com/ebuchman/go-shell-pipes

package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
)

func Execute(filename string, args ...string) (string, error) {
	return ExecuteEx(filename, "", false, args...)
}

var output string

func ExecuteEx(filename string, directoryName string, environment bool, args ...string) (string, error) {
	defer timeTrack(time.Now(), "execution")
	output = ""
	fullCommand := fmt.Sprintf("%s %s", filename, strings.Join(args[:], " "))
	log.Debugf("executing with %d args%s%s%sin directory %s with environment: %t", len(args), NewLine, fullCommand, NewLine, directoryName, environment)
	cmd := exec.Command(filename, args...)
	if directoryName != "" {
		cmd.Dir = directoryName
	}
	if environment {
		cmd.Env = os.Environ()
	}

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
	elapsed := time.Since(start)
	log.Debugf("%s took %ss", name, elapsed)
}

func writer(r io.Reader) {
	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output += string(buf[0:n])
			fmt.Print(string(buf[0:n]))
		}
		if err != nil {
			break
		}
	}
}
