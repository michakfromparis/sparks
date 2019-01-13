// FROM https://github.com/henrmota/errors-handling-example

package errx

import (
	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
)

func Warning(err error, message string) {
	if err != nil {
		log.Warnf("%s, cause: %s", message, err.Error())
	} else {
		log.Warn(message)
	}
}

func Error(err error, message string) {
	if err != nil {
		log.Errorf("%+v", errorx.Decorate(err, message))
	} else {
		log.Error(message)
	}
}

func Fatal(err error, message string) {
	if err != nil {
		log.Fatalf("%+v", errorx.Decorate(err, message))
	} else {
		log.Fatal(message)
	}
}