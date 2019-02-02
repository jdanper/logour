package main

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const appContext = "/v0"

// MountRoutes is a start point for all routes
func MountRoutes(r *fasthttprouter.Router) {

	r.GET(appContext+"/health", func(c *fasthttp.RequestCtx) {
		fmt.Fprintf(c, "OK")
	})

	r.OPTIONS(appContext+"/hit", func(c *fasthttp.RequestCtx) {
		c.Response.Header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Response.Header.Set("Access-Control-Allow-Origin", "*")
		c.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
		c.SetStatusCode(http.StatusNoContent)
	})

	r.POST(appContext+"/hit/:cid", func(ctx *fasthttp.RequestCtx) {
		xcred := ctx.Request.Header.Peek("X-CRED")
		reqInfo := extractReqInfo(ctx)
		info := make(chan *ClientInfo)

		go BuildHit(ctx.PostBody(), reqInfo, xcred, info)

		res := <-info
		header, err := json.Marshal(res)
		if err != nil {
			ctx.SetStatusCode(500)
			return
		}

		ctx.Response.Header.Add("X-CRED", base64.StdEncoding.EncodeToString(header))
		fmt.Fprint(ctx, "OK")
		ctx.SetStatusCode(200)

		ctx.SetConnectionClose()
	})
}

func extractReqInfo(ctx *fasthttp.RequestCtx) *RequestInfo {
	cid, _ := ctx.UserValue("cid").(string)
	return &RequestInfo{
		IP:        ctx.RemoteAddr().String(),
		Channel:   cid,
		UserAgent: string(ctx.UserAgent()),
	}
}
