package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_getSlackMessageFromEvent(t *testing.T) {
	e := &event{}
	_ = json.UnmarshalFromString(eventJSON, e)
	msg := getSlackMessageFromEvent(*e)

	assert.True(t, strings.Contains(msg, e.Message))
	assert.True(t, strings.Contains(msg, e.Client))
	assert.True(t, strings.Contains(msg, e.Kind))
}

func Test_mustNotify(t *testing.T) {
	fields := "ERROR,INFO,TRACE,DEBUG"
	assert.True(t, mustNotify("ERROR", fields), "should return true for valid event types")
	
	inexistantFields := "HELP,VALIDATION,SECURITY"
	assert.False(t, mustNotify("ERROR", inexistantFields))
}