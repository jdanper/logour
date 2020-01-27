package main

import (
	"log"
	"net"
	"time"

	_ "net/http/pprof"

	"bitbucket.org/danielper/util"
	cors "github.com/AdhityaRamadhanus/fasthttpcors"
	"github.com/buaazp/fasthttprouter"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

var (
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
	httpPort = util.GetEnvOrDefault("HTTP_PORT", "8080")
)

func main() {
	counter := time.Now()
	log.SetPrefix("[ LOGOUR ] ")

	db, err := connectScylla()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")
	defer db.Close()

	router, listener := setupHTTP(db)

	log.Println("Application running at :" + httpPort)
	log.Printf("Startup time: %s", time.Since(counter))

	corsHandler := getCorsHandler()

	log.Fatal(fasthttp.Serve(listener, corsHandler.CorsMiddleware(router.Handler)))
}

func getCorsHandler() *cors.CorsHandler {
	return cors.NewCorsHandler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: false,
		AllowMaxAge:      5600,
	})
}

func setupHTTP(db *DB) (*fasthttprouter.Router, net.Listener) {
	router := fasthttprouter.New()

	mountRoutes(router, db)

	ln, err := reuseport.Listen("tcp4", "0.0.0.0:"+httpPort)
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	return router, ln
}
