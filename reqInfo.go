package main

import (
	"strings"
)

var userAgentDevices = []string{
	"mobile",
	"windows",
	"cros",
	"macintosh",
	"x11",
}

const serverText = "server"

// RequestInfo represents the required values from request itself
type RequestInfo struct {
	IP        string
	UserAgent string
}

// GetDeviceType parses the user-agent string
func (req *RequestInfo) GetDeviceType() string {
	lowerUserAgnt := strings.ToLower(req.UserAgent)

	for _, v := range userAgentDevices {
		if strings.Contains(lowerUserAgnt, v) {
			return v
		}
	}

	return serverText
}
