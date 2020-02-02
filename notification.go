package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type SlackRequestBody struct {
	Text string `json:"text"`
}

var webHookUrl = os.Getenv("SLACK_HOOK_URL")
var plainRequiredKinds = os.Getenv("SLACK_NOTIFY_TYPES")

func notify(evt event) {
	if !mustNotify(evt.Kind, plainRequiredKinds) {
		return
	}

	slackJSONMessage := getSlackMessageFromEvent(evt)

	err := SendSlackNotification(webHookUrl, slackJSONMessage)
	if err != nil {
		log.Println("[NOTIFICATION] unable to send message to Slack")
		return
	}
}

func mustNotify(eventKind, plainRequiredKinds string) bool {
	requiredKinds := strings.ToLower(plainRequiredKinds)

	return strings.Contains(requiredKinds, strings.ToLower(eventKind))
}

func getSlackMessageFromEvent(evt event) string {
	msg := fmt.Sprintf(":bangbang: %s :bangbang: [%s] - [*%s*] - %s", evt.Kind, evt.CreatedAt.Format("2 Jan 2006 15:04:05"), evt.Client, evt.Message)

	if evt.JSONData != "" {
		msg = fmt.Sprintf("%s \n *Event data:* \n%s",msg, evt.JSONData)
	}

	return msg
}

func SendSlackNotification(webHookUrl string, msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webHookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	if buf.String() != "ok" {
		return errors.New("cannot send message to slack")
	}

	return nil
}
