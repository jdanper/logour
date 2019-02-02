package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/danielper/util"

	"github.com/buger/jsonparser"

	"bitbucket.org/danielper/util/msg"
)

const (
	hitPrefix     = "hit"
	sessionPrefix = "ssn"
	devicePrefix  = "dvc"
)

// RequestInfo represents the required values from request itself
type RequestInfo struct {
	IP        string
	Channel   string
	UserAgent string
}

type fullClientInfo struct {
	*Session
	*Device
}

// DvcType returns the type filtered from a user-agent string
func (req *RequestInfo) DvcType() string {
	dvcType := "DESKTOP"
	if strings.Contains(strings.ToLower(req.UserAgent), "mobile") {
		dvcType = "mobile"
	}

	return dvcType
}

// BuildHit creates a hit from
func BuildHit(event []byte, reqInfo *RequestInfo, xcred []byte, ids chan *ClientInfo) {
	info := make(chan *ClientInfo)

	session := getObject(event, "session")
	data := getObject(event, "data")
	dvcOtherIds := getObject(event, "deviceOtherIds")

	go FillCache(xcred, reqInfo, session, dvcOtherIds, info)

	cIds := <-info
	ids <- cIds

	compliance := checkCompliance(event)

	hit := jsonparser.StringToBytes("{}")

	hit, _ = jsonparser.Set(hit, compliance, "hit_iscompliable")
	hit, _ = jsonparser.Set(hit, toByteArray(reqInfo.IP), "ssn_ip")
	hit, _ = jsonparser.Set(hit, toByteArray(reqInfo.UserAgent), "ssn_useragent")
	hit, _ = jsonparser.Set(hit, toByteArray(reqInfo.Channel), "channel")
	hit, _ = jsonparser.Set(hit, toByteArray(reqInfo.DvcType()), "dvc_type")

	hit, _ = jsonparser.Set(hit, toByteArray(cIds.Sid.Value), "sid")
	hit, _ = jsonparser.Set(hit, toByteArray(cIds.Did.Value), "did")

	hit, _ = jsonparser.Set(hit, dvcOtherIds, "dvc_otherids")
	hit, _ = jsonparser.Set(hit, data, "hit_data")

	cpSession(hit, session)
	cpEventProps(event, hit)

	hit, _ = jsonparser.Set(hit, toByteArray(time.Now().Format("2006-01-02")), "hit_createdat")

	go produce(hit)
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

func cpSession(hit, session []byte) {
	mvProperty(session, hit, sessionPrefix, "referral")
	mvProperty(session, hit, sessionPrefix, "lat")
	mvProperty(session, hit, sessionPrefix, "lon")
	mvProperty(session, hit, sessionPrefix, "lang")
	mvProperty(session, hit, sessionPrefix, "apptype")
	mvProperty(session, hit, sessionPrefix, "appversion")
}

func cpEventProps(event, hit []byte) {
	mvProperty(event, hit, hitPrefix, "action")
	mvProperty(event, hit, hitPrefix, "context")
	mvProperty(event, hit, hitPrefix, "category")
	mvProperty(event, hit, hitPrefix, "label")

	mvProperty(event, hit, hitPrefix, "pageTitle")
	mvProperty(event, hit, hitPrefix, "lastScreen")
	mvProperty(event, hit, hitPrefix, "screen")
	mvProperty(event, hit, hitPrefix, "type")
	mvProperty(event, hit, hitPrefix, "windowWidth")
	mvProperty(event, hit, hitPrefix, "windowHeight")
}

func mvProperty(source, dest []byte, prefix, name string) {
	data, _, _, _ := jsonparser.Get(source, name)

	key := fmt.Sprintf("%s_%s", prefix, strings.ToLower(name))

	jsonparser.Set(source, data, key)
}

func checkCompliance(hit []byte) []byte {
	addtKeys := []string{"action", "label", "lastScreen", "screen", "data", "deviceOtherIds", "user", "session", "windowWidth", "windowHeight"}
	buff := &bytes.Buffer{}
	gob.NewEncoder(buff).Encode(addtKeys)
	additionalKeys := buff.Bytes()

	rK := []string{"channelId", "type", "context", "category", "pageTitle", "createdAt"}
	gob.NewEncoder(buff).Encode(rK)
	requiredKeys := buff.Bytes()

	allKeys := append(additionalKeys, requiredKeys...)

	var keys []byte
	jsonparser.EachKey(hit, func(indx int, key []byte, t jsonparser.ValueType, err error) {
		keys = append(keys, key...)
	})

	// required keys is part of keys
	if !bytes.Contains(keys, requiredKeys) {
		return []byte("false")
	}

	// keys is part of allkeys
	if !bytes.Contains(allKeys, keys) {
		return []byte("false")
	}

	if bytes.Contains(keys, []byte("user")) {
		uk := []string{"isLoggedIn", "uid", "email", "phone", "accountType", "primarySegment", "secondarySegment"}
		buff := &bytes.Buffer{}
		gob.NewEncoder(buff).Encode(uk)
		userKeys := buff.Bytes()

		user, _, _, err := jsonparser.Get(hit, "user")
		if err == nil {
			comp := "true"
			jsonparser.EachKey(user, func(ind int, key []byte, t jsonparser.ValueType, err error) {
				if !bytes.Contains(userKeys, key) {
					comp = "false"
				}
			})
			if comp == "false" {
				return []byte(comp)
			}
		}
	}

	if bytes.Contains(keys, []byte("session")) {
		sk := []string{"language", "os", "networkType", "lat", "lon", "osVersion", "appType", "appVersion"}
		buff := &bytes.Buffer{}
		gob.NewEncoder(buff).Encode(sk)
		userKeys := buff.Bytes()

		session, _, _, err := jsonparser.Get(hit, "session")
		if err == nil {
			comp := "true"
			jsonparser.EachKey(session, func(ind int, key []byte, t jsonparser.ValueType, err error) {
				if !bytes.Contains(userKeys, key) {
					comp = "false"
				}
			})

			if comp == "false" {
				return []byte(comp)
			}
		}
	}

	return []byte("true")
}

var topicPrefix = util.GetEnvOrDefault("KAFKA_TOPIC_PREFIX", "n13pqchi-")

func produce(content []byte) {
	channel, _ := jsonparser.GetString(content, "channel")
	// channel := topicPrefix + "TST8000"

	msg.Publish(content, topicPrefix+channel)
}
