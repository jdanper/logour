package main

import (
	"errors"
	"log"
	"time"

	"bitbucket.org/danielper/util/db"
	"github.com/globalsign/mgo/bson"
)

// NewKeySet returns a new ClientInfo object
func NewKeySet() *ClientInfo {
	return &ClientInfo{
		Sid: Sid{
			Value:     bson.NewObjectId().Hex(),
			ExpiresAt: time.Now().Add(time.Minute * 30).Local().String(),
		},
		Did: Did{Value: bson.NewObjectId().Hex()},
	}
}

// FillCache searches on cache store for given keys
func FillCache(clientInf []byte, reqInfo *RequestInfo, session, otherIds []byte, info chan *ClientInfo) {
	var fullClientInfo *fullClientInfo

	keys, err := ValidateKeys(clientInf)
	if err != nil {
		keys = NewKeySet()
		fullClientInfo = fullInfo(keys, reqInfo, otherIds)
	}

	info <- keys

	if fullClientInfo == nil {
		fullClientInfo, err = getIfExists(keys)
		if err != nil {
			fullClientInfo = fullInfo(keys, reqInfo, otherIds)
		}
	}
	go createCache(fullClientInfo)
	return
}

func getIfExists(data *ClientInfo) (*fullClientInfo, error) {
	device, err := FindDevice(data.Did.Value)
	if err != nil {
		return nil, errors.New("unable to verify did")
	}

	log.Println(device)

	session := &Session{}
	err = db.Mongo.C(SessionCollection).Find(bson.M{"_id": data.Sid.Value, "deviceId": data.Did.Value}).One(&session)
	if err != nil {
		log.Println("[Cache] Unable to verify sid.")
		return nil, errors.New("unable to verify did")
	}

	return &fullClientInfo{
		Device:  device,
		Session: session,
	}, nil
}

func fullInfo(info *ClientInfo, reqInfo *RequestInfo, otherIds []byte) *fullClientInfo {
	sess := &Session{
		ID:        bson.ObjectIdHex(info.Sid.Value),
		DeviceID:  bson.ObjectIdHex(info.Did.Value),
		IP:        reqInfo.IP,
		UserAgent: reqInfo.UserAgent,
		Channel:   reqInfo.Channel,
	}

	device := NewDevice(info.Did.Value, reqInfo.DvcType(), otherIds)

	return &fullClientInfo{
		Session: sess,
		Device:  device,
	}
}

func createCache(info *fullClientInfo) {
	err := info.Session.Create()
	if err != nil {
		log.Print(err.Error())
		log.Panicf("Unable to create session cache: %s", err.Error())
	}

	err = info.Device.Create()
	if err != nil {
		log.Panicln("Unable to create device cache")
	}
}
