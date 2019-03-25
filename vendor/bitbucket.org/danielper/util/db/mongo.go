package db

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"bitbucket.org/danielper/util"
	"github.com/globalsign/mgo"
)

// ConnectMongo returns a connection to a mongo database
func ConnectMongo(dbName string) (db *mgo.Database, err error) {
	mongoHost := util.GetEnvOrDefault("MONGO_HOST", "localhost")
	mongoUser := util.GetEnvOrDefault("MONGO_USER", "mongo")
	mongoPass := util.GetEnvOrDefault("MONGO_PASS", "mongo")

	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoUser, mongoPass, mongoHost, dbName)

	session := &mgo.Session{}

	if os.Getenv("GO_ENV") != "production" {
		session, err = mgo.Dial(mongoURL)
		if err != nil {
			return
		}
	} else {
		dialInfo, err := mgo.ParseURL(mongoURL)
		if err != nil {
			log.Fatal(err)
		}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			tlsConfig := &tls.Config{}
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			if err != nil {
				log.Fatal(err)
			}
			return conn, err
		}
		session, _ = mgo.DialWithInfo(dialInfo)
	}

	session.SetMode(mgo.Monotonic, true)

	db = session.DB(dbName)

	log.Printf("[Database] %s connected.", mongoURL)

	return
}
