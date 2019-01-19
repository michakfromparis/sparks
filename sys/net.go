package sys

import (
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// DownloadFile downloads the file at url and saves it to filepath
func DownloadFile(filepath string, url string) error {

	log.Debugf("sparks download %s %s", url, filepath)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
