package http

import (
	"io/ioutil"
	"net/http"
)

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Create the file
	return ioutil.WriteFile(filepath, dat, 0666)
}
