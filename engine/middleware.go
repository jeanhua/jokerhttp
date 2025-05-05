package engine

import (
	"encoding/json"
	"net/http"
)

type Middleware func(ctx *JokerContex)

type JokerContex struct {
	Request          *http.Request
	ResponseWriter   http.ResponseWriter
	MiddlewareChains []Middleware
	index            int
	maxIndex         int
	aborted          bool
}

func (ctx *JokerContex) Next() {
	if ctx.index < ctx.maxIndex && !ctx.aborted {
		ctx.index++
		ctx.MiddlewareChains[ctx.index](ctx)
	}
}

func (ctx *JokerContex) Abort() {
	ctx.aborted = true
}

func (ctx *JokerContex) AbortWithStatus(statusCode int) {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.Abort()
}

func (ctx *JokerContex) AbortWithStatusJSON(statusCode int, jsonObj interface{}) {
	jsonResult, err := json.Marshal(jsonObj)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.ResponseWriter.Write(jsonResult)
	ctx.Abort()
}

func (ctx *JokerContex) Use(middleware Middleware) {
	ctx.MiddlewareChains = append(ctx.MiddlewareChains, middleware)
	ctx.maxIndex = len(ctx.MiddlewareChains)
}
