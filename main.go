package main

import (
	"log"
	"net"
	"time"

	_ "net/http/pprof"

	"bitbucket.org/danielper/util/msg"
	"github.com/buaazp/fasthttprouter"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	counter := time.Now()
	log.SetPrefix("[ LOGOUR ] ")

	closeKafka := msg.NewKafkaProducer()

	defer closeKafka()

	router, listener := setupHTTPReusePort()

	log.Println("Application running at :8080")
	log.Println(time.Since(counter))

	log.Fatal(fasthttp.Serve(listener, router.Handler))
}

func setupHTTPReusePort() (*fasthttprouter.Router, net.Listener) {
	router := fasthttprouter.New()

	MountRoutes(router)

	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:8080", router.Handler))

	ln, err := reuseport.Listen("tcp4", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	return router, ln
}
