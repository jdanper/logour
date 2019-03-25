package db

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocql/gocql"
)

const (
	hostEnv = "CASSANDRA_HOST"
	userEnv = "CASSANDRA_USR"
	passEnv = "CASSANDRA_PWD"
)

// Cassandra points to active cassandra connection
var Cassandra *gocql.Session

// ConnectCassandra initiates a connection to Cassandra
func ConnectCassandra(keyspace string) (session *gocql.Session, err error) {
	plainHosts := os.Getenv(hostEnv)
	var hosts []string

	hosts = strings.Split(plainHosts, ",")

	if plainHosts == "" {
		hosts = []string{"127.0.0.1"}
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.ProtoVersion = 4
	cluster.Keyspace = keyspace

	username, password, err := getCredentials()
	if err != nil {
		log.Fatalln(err)
	}

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	session, err = cluster.CreateSession()
	Cassandra = session

	return
}

func getCredentials() (string, string, error) {
	username := os.Getenv(userEnv)
	password := os.Getenv(passEnv)
	var err error

	if strings.Contains(username, "/") {
		username, err = fromFile(username)
		if err != nil {
			return "", "", err
		}
	}
	if strings.Contains(password, "/") {
		password, err = fromFile(password)
		if err != nil {
			return "", "", err
		}
	}

	if username == "" {
		username = "cassandra"
	}
	if password == "" {
		password = "cassandra"
	}

	return username, password, nil
}

func fromFile(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(buff), nil
}
