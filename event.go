package main

import (
	"errors"
	"log"
	"time"

	"github.com/gocql/gocql"
)

var eventKinds = []string{"warning", "error", "info", "debug"}

type event struct {
	ID         gocql.UUID
	Hostname   string              `json:"hostname"`
	Kind       string              `json:"type"`
	Action     string              `json:"action"`
	Message    string              `json:"message"`
	CreatedAt  time.Time           `json:"created_at"`
	CustomData map[string]struct{} `json:"custom"`

	IP        string
	UserAgent string
	SavedAt   time.Time
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

	content := &event{}
	json.Unmarshal(rawEvent, content)

	if !content.isValid() {
		return nil, errors.New("invalid event data")
	}

	content.IP = reqInfo.IP
	content.UserAgent = reqInfo.UserAgent
	content.SavedAt = time.Now()

	return content, nil
}

func (e *event) isValid() bool {
	valid := (e.Hostname != "" &&
		e.Action != "" &&
		e.Message != "")

	valid = e.Kind != "" && contains(e.Kind)

	return valid
}

func contains(src string) bool {
	for _, n := range eventKinds {
		if src == n {
			return true
		}
	}
	return false
}
