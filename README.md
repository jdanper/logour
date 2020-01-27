# Logour
[![Go Report Card](https://goreportcard.com/badge/github.com/jdanper/logour)](https://goreportcard.com/report/github.com/jdanper/logour)

Http log collector backed by Apache Cassandra/ScyllaDB

### Building

** Make sure [dep](https://github.com/golang/dep/blob/master/docs/installation.md) is installed

Then build with `make build`. 
Under the hood, this command will clean, install dependencies, test and build the binary

### How to test

Just run `make test`
