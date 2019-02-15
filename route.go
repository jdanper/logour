package main

import (
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const appContext = "/v0"

// MountRoutes is a start point for all routes
func MountRoutes(r *fasthttprouter.Router) {

	r.GET(appContext+"/health", func(c *fasthttp.RequestCtx) {
		fmt.Fprintf(c, "OK")
	})

	r.POST(appContext+"/event", func(ctx *fasthttp.RequestCtx) {
		reqInfo := extractReqInfo(ctx)

		go BuildEvent(ctx.PostBody(), reqInfo)

		fmt.Fprint(ctx, "OK")

		ctx.SetStatusCode(200)

		ctx.SetConnectionClose()
	})
}

func extractReqInfo(ctx *fasthttp.RequestCtx) *RequestInfo {
	client := ctx.Request.Header.Peek("X-LOGOUR-TOKEN")

	return &RequestInfo{
		IP:        ctx.RemoteAddr().String(),
		Client:    string(client),
		UserAgent: string(ctx.UserAgent()),
	}
}
