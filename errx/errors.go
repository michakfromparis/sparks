// FROM https://github.com/henrmota/errors-handling-example

package errx

import (
	log "github.com/sirupsen/logrus"
	"github.com/joomcode/errorx"
)

// Fatal is used terminate the program with an error
func Fatal(err error) {
	log.Fatalf("%v", err)
}

// Fatalf is used terminate the program with an error and decorate the last error with a message
func Fatalf(err error, message string) {
	if err != nil {
		Fatal(errorx.Decorate(err, message))
	} else {
		log.Fatalf(message)
	}
}
