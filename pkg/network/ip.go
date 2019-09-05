package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"time"
)

const regexIP = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`

var searchIP = regexp.MustCompile(regexIP).FindString

// Get preferred outbound ip of this machine
func GetOutboundIP() (IP string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func GetPublicIP() (IP string, err error) {
	httpClient := &http.Client{Timeout: time.Second}
	response, err := httpClient.Get("https://duckduckgo.com?q=ip")
	if err != nil {
		return "", err
	} else if response.StatusCode != 200 {
		return "", fmt.Errorf("https://duckduckgo.com?q=ip failed with status code %s", response.Status)
	}
	content, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return "", err
	}
	IP = searchIP(string(content))
	if IP == "" {
		return "", fmt.Errorf("Public IP not found")
	}
	return IP, nil
}
