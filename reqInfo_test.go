package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var reqInfo = &RequestInfo{
	IP:        "127.0.0.1",
	UserAgent: "test user agent",
}

func TestRequestInfo_GetDeviceType(t *testing.T) {
	device := reqInfo.GetDeviceType()

	assert.Equal(t, "server", device, "resolved user agent must be server")
}

func TestRequestInfo_GetDesktopDeviceType(t *testing.T) {
	reqInfo.UserAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"
	device := reqInfo.GetDeviceType()

	assert.Equal(t, "windows", device, "resolved user agent must be server")
}

func TestRequestInfo_GetMobileDeviceType(t *testing.T) {
	reqInfo.UserAgent = "Mozilla/5.0 (Linux; Android 7.0; SM-G892A Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/60.0.3112.107 Mobile Safari/537.36"
	device := reqInfo.GetDeviceType()

	assert.Equal(t, "mobile", device, "resolved user agent must be server")
}
