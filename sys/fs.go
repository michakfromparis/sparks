package sys

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/michakfromparis/sparks/errx"
)

// MkDir creates the directory at path and all subdirectories
func MkDir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		errx.Fatalf(err, "failed to create directory: "+path)
	}
	return nil
}

// DirSize returns a human readable form of the size of the directory at path
func DirSize(path string) (string, error) {
	du, err := Execute("du", "-hs", path)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(du)
	if len(fields) > 0 {
		return fields[0], nil
	}
	return "", nil
}

// Pwd returns the current directory's full path
func Pwd() (string, error) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		errx.Fatalf(err, "Could not determine current working directory")
	}
	log.Debug("pwd: " + path)
	return path, err
}

// SedFile searches regex inside filename and replaces it by newContent
func SedFile(filename string, regex string, newContent string) {
	log.Tracef("sed %s %s in %s", regex, newContent, filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errx.Fatalf(err, "Could not read: "+filename)
	}
	re := regexp.MustCompile(regex)
	// log.Trace(string(bytes))
	matched := re.MatchString(string(bytes))
	if !matched {
		log.Warnf("regex %s did not match content of %s", regex, filename)
	}
	data := re.ReplaceAllString(string(bytes), newContent)
	// log.Trace(data)
	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		errx.Fatalf(err, "Could not write to: "+filename)
	}
}
