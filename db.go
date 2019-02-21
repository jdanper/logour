package main

import (
	"log"
	"strings"

	"bitbucket.org/danielper/util"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

var (
	hosts    = util.GetEnvOrDefault("SCYLLA_HOSTS", "localhost")
	keyspace = util.GetEnvOrDefault("SCYLLA_KEYSPACE", "logour")

	session *gocql.Session

	columns = []string{"id", "hostname", "kind", "action", "message", "custom_data", "ip", "user_agent", "saved_at", "created_at"}
)

func connectScylla() (*gocql.Session, error) {
	cluster := gocql.NewCluster(strings.Split(hosts, ",")...)

	cluster.Keyspace = keyspace
	// cluster.ProtoVersion = 4

	cluster.Consistency = gocql.Any

	var err error
	session, err = cluster.CreateSession()

	return session, err
}

func insert(content *event) {
	content.ID = gocql.TimeUUID()

	stmt, names := qb.Insert(keyspace + ".event").Columns(columns...).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(content)

	if err := q.ExecRelease(); err != nil {
		log.Println(err)
	}
}
