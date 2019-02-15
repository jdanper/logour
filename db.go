package main

import (
	"strings"

	"bitbucket.org/danielper/util"

	"github.com/gocql/gocql"
)

func connectScylla() (*gocql.Session, error) {
	hosts := util.GetEnvOrDefault("SCYLLA_HOSTS", "localhost")
	keyspace := util.GetEnvOrDefault("SCYLLA_KEYSPACE", "logour")

	cluster := gocql.NewCluster(strings.Split(hosts, ",")...)

	cluster.Keyspace = keyspace

	cluster.Consistency = gocql.Quorum

	return cluster.CreateSession()
}

func insert() {

}
