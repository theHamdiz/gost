package dwn

import (
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	return err
}
