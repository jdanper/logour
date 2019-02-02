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

// Mongo contains an active connection to MongoDB
var Mongo *mgo.Database

// ConnectMongo returns a connection to a mongo database
func ConnectMongo(dbName string) (session *mgo.Session, err error) {
	mongoHosts := util.GetEnvOrDefault("MONGO_HOSTS", "localhost:27017")
	mongoUser := util.GetEnvOrDefault("MONGO_USERNAME", "mongo")
	mongoPassword := util.GetEnvOrDefault("MONGO_PASSWORD", "mongo")

	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoUser, mongoPassword, mongoHosts, dbName)

	session = &mgo.Session{}

	if os.Getenv("GO_ENV") != "production" {
		session, err = mgo.Dial(mongoURL)
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

	if err != nil {
		return
	}

	session.SetMode(mgo.Monotonic, true)

	Mongo = session.DB(os.Getenv("db"))

	log.Printf("[Database] %s connected.", mongoURL)

	return
}
