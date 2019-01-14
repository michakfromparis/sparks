package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/errx"
)

func Pwd() (string, error) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		errx.Fatalf(err, "Could not determine current working directory")
	}
	log.Debug("pwd: " + path)
	return path, err
}

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
