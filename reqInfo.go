package main

import (
	"strings"
)

const (
	mobileName  = "mobile"
	desktopName = "desktop"
	serverName  = "server"
)

// RequestInfo represents the required values from request itself
type RequestInfo struct {
	IP        string
	Client    string
	UserAgent string
}

// GetDeviceType parses the user-agent string
func (req *RequestInfo) GetDeviceType() string {

	dvcType := serverName

	lowerUserAgnt := strings.ToLower(req.UserAgent)

	if strings.Contains(lowerUserAgnt, desktopName) {
		dvcType = desktopName
	}

	if strings.Contains(lowerUserAgnt, mobileName) {
		dvcType = mobileName
	}

	return dvcType
}
