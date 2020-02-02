package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

var eventKinds = []string{"WARNING", "ERROR", "INFO", "DEBUG"}

type payload struct {
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
	*payload
	CreatedAt time.Time
	JSONData  string
}

func createEvent(rawEvent []byte, reqInfo *RequestInfo, db Database) error {
	evt, err := buildEvent(rawEvent, reqInfo)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	go db.Insert(evt)
	go notify(*evt)

	return nil
}

func buildEvent(rawEvent []byte, reqInfo *RequestInfo) (*event, error) {
	event, err := parseEvent(rawEvent)
	if err != nil {
		return nil, err
	}

	if err = event.checkValid(); err != nil {
		return nil, err
	}

	event.payload.RemoteAddress = reqInfo.IP
	event.payload.UserAgent = reqInfo.UserAgent
	event.payload.SavedAt = time.Now()

	createdAt := time.Unix(event.payload.CreatedAt, 0)
	event.CreatedAt = createdAt

	customData, _ := json.Marshal(event.payload.CustomJSON)

	event.JSONData = string(customData)

	return event, nil
}

func parseEvent(rawEvent []byte) (*event, error) {
	if !json.Valid(rawEvent) {
		return nil, errInvalidJSON
	}

	eventPayload := &payload{}
	_ = json.Unmarshal(rawEvent, eventPayload)

	return &event{payload: eventPayload}, nil
}

func (e *payload) checkValid() error {
	containsMandatory := e.Hostname != "" && e.Message != "" && e.Client != ""

	if !containsMandatory {
		return errEmptyMandatoryField
	}

	if !(e.Kind != "" && containsKind(e.Kind)) {
		return errInvalidKind
	}

	return nil
}

func containsKind(kind string) bool {
	for _, n := range eventKinds {
		if kind == n {
			return true
		}
	}

	return false
}

func (e *event) String() string {
	return fmt.Sprintf("%s - %s - %s - %s", e.CreatedAt, e.Kind, e.Client, e.Message)
}