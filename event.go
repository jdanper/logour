package main

import (
	"errors"
	"log"
	"time"

	"github.com/gocql/gocql"
)

var eventKinds = []string{"WARNING", "ERROR", "INFO", "DEBUG"}

type eventPayload struct {
	ID         gocql.UUID
	Client     string                 `json:"client"`
	Hostname   string                 `json:"hostname"`
	Kind       string                 `json:"type"`
	Message    string                 `json:"message"`
	CreatedAt  int64                  `json:"createdAt"`
	CustomJSON map[string]interface{} `json:"custom"`

	RemoteAddress string
	UserAgent     string
	SavedAt       time.Time
}

type event struct {
	*eventPayload
	CreatedAt time.Time
	JSONData  string
}

func process(rawEvent []byte, reqInfo *RequestInfo) {
	evt, err := buildEvent(rawEvent, reqInfo)

	if err != nil {
		log.Println(err.Error())
		return
	}

	go insert(evt)
}

func buildEvent(rawEvent []byte, reqInfo *RequestInfo) (*event, error) {
	if !json.Valid(rawEvent) {
		return nil, errors.New("invalid json payload")
	}

	content := &eventPayload{}
	json.Unmarshal(rawEvent, content)

	if !content.isValid() {
		return nil, errors.New("invalid event data")
	}

	event := &event{eventPayload: content}

	content.RemoteAddress = reqInfo.IP
	content.UserAgent = reqInfo.UserAgent
	content.SavedAt = time.Now()

	createdAt := time.Unix(content.CreatedAt, 0)
	event.CreatedAt = createdAt

	custJSON, err := json.Marshal(content.CustomJSON)
	if err != nil {
		log.Println("unable to parse map to json string")
	}

	event.JSONData = string(custJSON)

	log.Println(event.JSONData)

	log.Println(custJSON)

	return event, nil
}

func (e *eventPayload) isValid() bool {
	valid := (e.Hostname != "" && e.Message != "" && e.Client != "")

	valid = e.Kind != "" && containsKind(e.Kind)

	return valid
}

func containsKind(src string) bool {
	for _, n := range eventKinds {
		if src == n {
			return true
		}
	}

	return false
}
