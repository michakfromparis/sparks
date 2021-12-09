package sys

import (
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// DownloadFile downloads the file at url and saves it to filepath
func DownloadFile(filepath string, url string) (rerr error) {

	log.Debugf("sparks download %s %s", url, filepath)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			rerr = err
		}
	}()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			rerr = err
		}
	}()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
