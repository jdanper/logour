package main

import (
	"bitbucket.org/danielper/util/db"
	"github.com/globalsign/mgo/bson"
)

// Session contains the info about user session on client channel
type Session struct {
	ID        bson.ObjectId
	DeviceID  bson.ObjectId
	IP        string
	UserAgent string
	Channel   string
}

// SessionCollection indicates the database collection of the session entity
const SessionCollection = "sessions"

// Create inserts a new session on database
func (s *Session) Create() error {
	return db.Mongo.C(SessionCollection).Insert(s)
}

// FindSession searches for a session with given sid in database
func FindSession(sid string) (s *Session, err error) {
	err = db.Mongo.C(SessionCollection).FindId(sid).One(&s)
	return
}
