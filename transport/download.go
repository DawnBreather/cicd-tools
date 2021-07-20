package transport

import (
	"github.com/DawnBreather/go-commons/logger"
	"io"
	"net/http"
	"os"
)

var _logger = logger.New()

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		_logger.Errorf("Unable to download { %s }: %v", url, err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		_logger.Errorf("Unable to create file { %s }: %v", url, filepath)
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		_logger.Errorf("Unable to save downloaded file { %s } into { %s }: %v", url, filepath)
		return
	}
}