package network

import (
	"fmt"
	"net/http"

	netlib "github.com/qdm12/golibs/network"
)

func (n *network) CheckMultipleHTTPConnections(urls []string) (errors []error) {
	chErr := make(chan error)
	for i := range urls {
		go func(url string) {
			chErr <- checkHTTPConnection(n.client, url)
		}(urls[i])
	}
	for range urls {
		err := <-chErr
		if err != nil {
			errors = append(errors, err)
		}
	}
	close(chErr)
	return errors
}

func checkHTTPConnection(client netlib.Client, url string) (err error) {
	_, status, err := client.GetContent(url)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("connectivity to %s failed: HTTP status code is %d", url, status)
	}
	return nil
}
