package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type response struct {
	Message string `json:"message"`
}

const appContext = "/logour/v1"

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func mountRoutes(r *fasthttprouter.Router, db *DB) {

	r.GET(appContext+"/health", func(c *fasthttp.RequestCtx) {
		_, err := fmt.Fprintf(c, "OK")
		if err != nil {
			log.Println("Unable to saveEvent http request")
		}
	})

	r.POST(appContext+"/event", func(ctx *fasthttp.RequestCtx) {
		if len(ctx.PostBody()) == 0 {
			ctx.Response.SetStatusCode(400)
			return
		}

		reqInfo := extractReqInfo(ctx)

		go createEvent(ctx.PostBody(), reqInfo, db)

		end(ctx, 200)
	})
}

func extractReqInfo(ctx *fasthttp.RequestCtx) *RequestInfo {
	return &RequestInfo{
		IP:        ctx.RemoteAddr().String(),
		UserAgent: string(ctx.UserAgent()),
	}
}

func end(ctx *fasthttp.RequestCtx, code int) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)

	ctx.SetBody([]byte("{\"message\":\"OK\"}"))

	ctx.SetConnectionClose()
}
