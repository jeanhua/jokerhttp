package engine

import (
	"encoding/json"
	"net/http"
)

type Middleware func(ctx *Contex)

type Contex struct {
	Request          *http.Request
	ResponseWriter   http.ResponseWriter
	MiddlewareChains []Middleware
	index            int
	maxIndex         int
	aborted          bool
}

func (ctx *Contex) Next() {
	if ctx.index < ctx.maxIndex && !ctx.aborted {
		ctx.index++
		ctx.MiddlewareChains[ctx.index](ctx)
	}
}

func (ctx *Contex) Abort() {
	ctx.aborted = true
}

func (ctx *Contex) AbortWithStatus(statusCode int) {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.Abort()
}

func (ctx *Contex) AbortWithStatusJSON(statusCode int, jsonObj interface{}) {
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

func (ctx *Contex) Use(middleware Middleware) {
	ctx.MiddlewareChains = append(ctx.MiddlewareChains, middleware)
	ctx.maxIndex = len(ctx.MiddlewareChains)
}
