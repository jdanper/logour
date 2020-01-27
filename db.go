package main

import (
	"log"
	"strings"

	"bitbucket.org/danielper/util"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

// DB wraps a database connection
type DB struct {
	session *gocql.Session
}

// Database defines basic database operations
type Database interface {
	Close()
	Insert(evt *event)
}

var (
	hosts    = util.GetEnvOrDefault("SCYLLA_HOSTS", "localhost")
	keyspace = util.GetEnvOrDefault("SCYLLA_KEYSPACE", "logour")

	columns = []string{"id", "client", "hostname", "kind", "message", "json_data", "remote_address", "user_agent", "saved_at", "created_at"}
)

func connectScylla() (*DB, error) {
	cluster := gocql.NewCluster(strings.Split(hosts, ",")...)

	cluster.Keyspace = keyspace

	cluster.Consistency = gocql.Any

	session, err := cluster.CreateSession()

	return &DB{session}, err
}

// Insert persists an event
func (db *DB) Insert(content *event) {
	content.ID = gocql.TimeUUID()

	stmt, names := qb.Insert(keyspace + ".event").Columns(columns...).ToCql()

	q := gocqlx.Query(db.session.Query(stmt), names).BindStruct(content)

	if err := q.ExecRelease(); err != nil {
		log.Println(err)
	}
}

// Close ends database connection
func (db *DB) Close() {
	db.session.Close()
}
