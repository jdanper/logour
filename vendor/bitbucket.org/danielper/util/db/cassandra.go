package db

import (
	"os"

	"github.com/gocql/gocql"
)

// Cassandra points to active cassandra connection
var Cassandra *gocql.Session

// ConnectCassandra initiates a connection to Cassandra
func ConnectCassandra() (err error) {
	host := os.Getenv("CASSANDRA_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	cluster := gocql.NewCluster(host)
	cluster.ProtoVersion = 4
	cluster.Keyspace = "iicdata"

	username := os.Getenv("CASSANDRA_USR")
	password := os.Getenv("CASSANDRA_PWD")
	if username == "" {
		username = "cassandra"
	}
	if password == "" {
		password = "cassandra"
	}

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	Cassandra, err = cluster.CreateSession()

	return
}
