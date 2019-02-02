package main

import (
	"encoding/base64"
)

// Sid represents a session id used by client
type Sid struct {
	Value     string `json:"value"`
	ExpiresAt string `json:"expiresAt"`
}

// Did represents a device id used by client
type Did struct {
	Value string `json:"value"`
}

// ClientInfo represents a set of key used to identify user devices and its sessions
type ClientInfo struct {
	Sid Sid `json:"sid"`
	Did Did `json:"did"`
}

// ValidateKeys makes sure the xCred passed is valid
func ValidateKeys(xCred []byte) (*ClientInfo, error) {
	if len(xCred) == 0 {
		return NewKeySet(), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(string(xCred))
	if err != nil {
		return nil, err
	}

	data := &ClientInfo{}
	err = json.Unmarshal([]byte(decoded), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// EncodeKeys returns a base64 encoded string
func EncodeKeys(xcred []byte) string {
	return base64.StdEncoding.EncodeToString(xcred)
}
