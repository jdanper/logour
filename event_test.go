package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const eventJSON = `
{
	"client": "sample app",
	"hostname": "sample host",
	"type": "ERROR",
	"message":"A test error",
	"createdAt": 1579914374,
	"custom":{}
}`

const invalidEventJSON = `
{
	"client": "sample app",
	"hostname": "sample host",
	"type": "some type",
	"message":"A test error",
	"createdAt": 1579914374,
	"custom":{}
}`

const malformedEventJSON = `
{
	"client: "sample app",
	"hostname": "sample host",
	"type": "some type",
	"message":"A test error",
	"createdAt": 1579914374,
	"custom":
}`

const malformedCustomEventJSON = `
{
	"client": "sample app",
	"hostname": "sample host",
	"type": "some type",
	"message":"A test error",
	"createdAt": 1579914374,
	"custom":{--}
}`

var requestInfo = &RequestInfo{
	IP:        "127.0.0.1",
	UserAgent: "test agent",
}

func BenchmarkBuildEvent(b *testing.B) {
}

func Test_buildEventWithoutError(t *testing.T) {
	event, err := buildEvent([]byte(eventJSON), requestInfo)

	assert.NoError(t, err, "build event should not return error")
	assert.NotNil(t, event, "event must be valid")
}

func Test_buildEventWithInvalidEvent(t *testing.T) {
	event, err := buildEvent([]byte(invalidEventJSON), requestInfo)

	assert.Error(t, err, "build event should return error")
	assert.Nil(t, event, "event is not valid")
}

func Test_buildEventWithMalformedPayload(t *testing.T) {
	event, err := buildEvent([]byte(malformedEventJSON), requestInfo)

	assert.Error(t, err, "build event should return error")
	assert.Nil(t, event, "event is not valid")
}

func Test_containsKind(t *testing.T) {
	containsAllValidKinds := false

	for _, v := range eventKinds {
		containsAllValidKinds = containsKind(v)
	}

	assert.True(t, containsAllValidKinds, "must pass all valid event kinds")
}

func Test_containsInvalidKind(t *testing.T) {
	assert.False(t, containsKind("test"), "must return false to invalid kinds")
}

func Test_parseValidEvent(t *testing.T) {
	event, err := parseEvent([]byte(eventJSON))

	assert.NoError(t, err, "should parse the event successfully")
	assert.NotNil(t, event, "parsed event must be present")
}

func Test_parseInvalidEvent(t *testing.T) {
	event, err := parseEvent([]byte(malformedEventJSON))

	assert.Error(t, err, "should not parse the event")
	assert.Nil(t, event, "must return nil event")
}

func Test_payload_checkValid(t *testing.T) {
	event, _ := parseEvent([]byte(eventJSON))

	err := event.payload.checkValid()

	assert.Nil(t, err, "payload must be treated as valid")
}

func Test_notContainsMandatoryFields(t *testing.T) {
	event, _ := parseEvent([]byte(eventJSON))

	event.Client = ""
	err := event.payload.checkValid()

	assert.Error(t, err, "should block invalid or empty mandatory fields")
}

type testDB struct{}

func (db *testDB) Insert(content *event) {}
func (db *testDB) Close()                {}

func Test_createEvent(t *testing.T) {
	err := createEvent([]byte(eventJSON), requestInfo, &testDB{})

	assert.NoError(t, err, "must create without errors")
}

func Test_createInvalidEvent(t *testing.T) {
	err := createEvent([]byte(invalidEventJSON), requestInfo, &testDB{})

	assert.Error(t, err, "must end with error")
}
