package network

import (
	"fmt"
	"net/http"
	"time"
)

func CheckMultipleHTTPConnections(URLs []string) (errors []error) {
	httpClient := &http.Client{Timeout: time.Second}
	chErr := make(chan error)
	for i := range URLs {
		URL := URLs[i]
		go func() {
			chErr <- checkHTTPConnection(httpClient, URL)
		}()
	}
	N := len(URLs)
	for N > 0 {
		select {
		case err := <-chErr:
			if err != nil {
				errors = append(errors, err)
			}
			N--
		}
	}
	close(chErr)
	return errors
}

func checkHTTPConnection(httpClient *http.Client, URL string) (err error) {
	response, err := httpClient.Get(URL)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return fmt.Errorf("connectivity to %s failed: HTTP status code is %s", URL, response.Status)
	}
	return nil
}
