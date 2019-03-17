package main

import (
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type response struct {
	Message string `json:"message"`
}

const appContext = "/logour"

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func mountRoutes(r *fasthttprouter.Router) {

	r.GET(appContext+"/health", func(c *fasthttp.RequestCtx) {
		fmt.Fprintf(c, "OK")
	})

	r.POST(appContext+"/event", func(ctx *fasthttp.RequestCtx) {
		if len(ctx.PostBody()) == 0 {
			ctx.Response.SetStatusCode(400)
			return
		}

		reqInfo := extractReqInfo(ctx)

		go process(ctx.PostBody(), reqInfo)

		end(ctx, 200, &response{Message: "OK"})
	})
}

func extractReqInfo(ctx *fasthttp.RequestCtx) *RequestInfo {
	return &RequestInfo{
		IP:        ctx.RemoteAddr().String(),
		UserAgent: string(ctx.UserAgent()),
	}
}

func end(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)

	ctx.SetBody([]byte("{\"message\":\"OK\"}"))

	ctx.SetConnectionClose()
}
