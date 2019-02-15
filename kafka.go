package main

import (
	"bitbucket.org/danielper/util"
	"bitbucket.org/danielper/util/msg"
)

var topic = util.GetEnvOrDefault("KAFKA_TOPIC", "n13pqchi-event")

func produce(content []byte) {
	msg.Publish(content, topic)
}
