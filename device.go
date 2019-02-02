package main

import (
	"time"

	"bitbucket.org/danielper/util/db"

	"github.com/globalsign/mgo/bson"
)

// Device represents a client device
type Device struct {
	ID        bson.ObjectId
	Type      string
	OtherIds  []byte
	CreatedAt time.Time
}

// DeviceCollection is the collection name in mongodb
const DeviceCollection = "devices"

// Create inserts a new Device into the collection
func (d *Device) Create() error {
	return db.Mongo.C(DeviceCollection).Insert(d)
}

// FindDevice retrieves a Device with given id or nil if none is found
func FindDevice(did string) (d *Device, err error) {
	err = db.Mongo.C(DeviceCollection).FindId(did).One(&d)
	return
}

// NewDevice creates a new device struct with given content
func NewDevice(id string, dvcType string, otherIds []byte) *Device {
	return &Device{
		ID:        bson.ObjectIdHex(id),
		Type:      dvcType,
		OtherIds:  otherIds,
		CreatedAt: time.Now(),
	}
}
