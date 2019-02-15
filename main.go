package main

import (
	"log"
	"net"
	"time"

	_ "net/http/pprof"

	cors "github.com/AdhityaRamadhanus/fasthttpcors"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

func main() {
	counter := time.Now()
	log.SetPrefix("[ LOGOUR ] ")

	dbSession, err := connectScylla()
	if err != nil {
		log.Fatal(err)
	}

	defer dbSession.Close()

	router, listener := setupHTTPReusePort()

	log.Println("Application running at :8080")
	log.Println(time.Since(counter))

	corsHandler := getCorsHandler()

	log.Fatal(fasthttp.Serve(listener, corsHandler.CorsMiddleware(router.Handler)))
}

func getCorsHandler() *cors.CorsHandler {
	return cors.NewCorsHandler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"x-something-client", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: false,
		AllowMaxAge:      5600,
	})
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
