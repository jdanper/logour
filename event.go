package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/buger/jsonparser"
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

// BuildEvent creates a event from
func BuildEvent(rawEvent []byte, reqInfo *RequestInfo) {
	if !isEventCompliant(rawEvent) {
		return
	}

	event := jsonparser.StringToBytes("{}")

	event = extractBasicData(rawEvent, reqInfo)
	custom := getObject(rawEvent, "data")

	event, _ = jsonparser.Set(event, toByteArray(reqInfo.IP), "ip")
	event, _ = jsonparser.Set(event, toByteArray(reqInfo.UserAgent), "useragent")
	event, _ = jsonparser.Set(event, toByteArray(reqInfo.Client), "client_id")
	event, _ = jsonparser.Set(event, toByteArray(reqInfo.GetDeviceType()), "dvc_type")

	event, _ = jsonparser.Set(event, custom, "custom_data")

	event, _ = jsonparser.Set(event, toByteArray(time.Now().Format("2006-01-02")), "event_savedat")

	go produce(event)
}

func extractBasicData(event []byte, reqInfo *RequestInfo) []byte {
	basicData, _, _, _ := jsonparser.Get(event, "hostname")
	basicData, _, _, _ = jsonparser.Get(event, "type")
	basicData, _, _, _ = jsonparser.Get(event, "action")
	basicData, _, _, _ = jsonparser.Get(event, "message")
	basicData, _, _, _ = jsonparser.Get(event, "user")

	return basicData
}

func toByteArray(content string) []byte {
	return []byte(fmt.Sprintf(`"%s"`, content))
}

func getObject(source []byte, key string) []byte {
	obj, _, _, err := jsonparser.Get(source, key)

	if err != nil {
		return jsonparser.StringToBytes("{}")
	}

	return obj
}

func isEventCompliant(event []byte) bool {
	buff := &bytes.Buffer{}

	requiredKeyNames := []string{"hostname", "type", "action", "message"}
	gob.NewEncoder(buff).Encode(requiredKeyNames)

	requiredKeys := buff.Bytes()

	var keys []byte
	jsonparser.EachKey(event, func(indx int, key []byte, t jsonparser.ValueType, err error) {
		keys = append(keys, key...)
	})

	if !bytes.Contains(keys, requiredKeys) {
		return false
	}

	return true
}
