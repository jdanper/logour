# Logour
[![Go Report Card](https://goreportcard.com/badge/github.com/jdanper/logour)](https://goreportcard.com/report/github.com/jdanper/logour)

Http log collector backed by Apache Cassandra/ScyllaDB

### Building

** Make sure [dep](https://github.com/golang/dep/blob/master/docs/installation.md) is installed

Then build with `make build`. 
Under the hood, this command will clean, install dependencies, test and build the binary

### Database

You have to setup the database with the scripts from `/scripts` and expose port `9042`

### How to test

Just run `make test`

### How to Run

Make sure database is up and running

Environment variables needed:
* `HTTP_PORT` defaults to `8080`
* `DATABASE_HOSTS` defaults to `localhost`
* `DATABASE_KEYSPACE` defaults to `logour`

Build the binary then run with `./logour`