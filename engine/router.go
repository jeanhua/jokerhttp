package engine

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type JokerRouter struct {
	prefix      string
	engine      *JokerEngine
	middlewares []Middleware
}

func (engine *JokerEngine) NewRouter() *JokerRouter {
	return &JokerRouter{
		prefix: "",
		engine: engine,
	}
}

func (router *JokerRouter) Group(prefix string) *JokerRouter {
	if !strings.HasPrefix(prefix, "/") {
		panic("[Error]:Prefix must start with /")
	}
	if router.prefix == "/" {
		return &JokerRouter{
			prefix:      prefix,
			engine:      router.engine,
			middlewares: router.middlewares,
		}
	} else {
		return &JokerRouter{
			prefix:      router.prefix + prefix,
			engine:      router.engine,
			middlewares: router.middlewares,
		}
	}
}

func (router *JokerRouter) Use(middleware Middleware) {
	router.middlewares = append(router.middlewares, middleware)
}

func (router *JokerRouter) Map(pattern string, handle func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{})) {
	pattern = router.prefix + pattern
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount_global := len(router.engine.middlewares)
		middlewareCount_router := len(router.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount_global+middlewareCount_router+1),
			index:            -1,
			maxIndex:         middlewareCount_global + middlewareCount_router + 1,
		}
		copy(ctx.MiddlewareChains, router.engine.middlewares)
		copy(ctx.MiddlewareChains[middlewareCount_global:], router.middlewares)
		finalHandler := func(ctx *JokerContex) {
			params := r.URL.Query()
			status, response := handle(r, params, func(key, value string) {
				w.Header().Set(key, value)
			})
			if response == nil {
				w.WriteHeader(status)
				return
			}
			jsonresult, err := json.Marshal(response)
			if err != nil {
				log.Println("[Error]:Handle in " + pattern + " >>> " + err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Server", "JokerHttp")
			w.WriteHeader(status)
			w.Write(jsonresult)
		}
		ctx.MiddlewareChains[middlewareCount_global+middlewareCount_router] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func (router *JokerRouter) MapGet(pattern string, handle func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{})) {
	pattern = router.prefix + pattern
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount_global := len(router.engine.middlewares)
		middlewareCount_router := len(router.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount_global+middlewareCount_router+1),
			index:            -1,
			maxIndex:         middlewareCount_global + middlewareCount_router + 1,
		}
		copy(ctx.MiddlewareChains, router.engine.middlewares)
		copy(ctx.MiddlewareChains[middlewareCount_global:], router.middlewares)
		finalHandler := func(ctx *JokerContex) {
			if r.Method != http.MethodGet {
				w.WriteHeader(405)
				return
			}
			params := r.URL.Query()
			status, response := handle(r, params, func(key, value string) {
				w.Header().Set(key, value)
			})
			if response == nil {
				w.WriteHeader(status)
				return
			}
			jsonresult, err := json.Marshal(response)
			if err != nil {
				log.Println("[Error]:Handle in " + pattern + " >>> " + err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Server", "JokerHttp")
			w.WriteHeader(status)
			w.Write(jsonresult)
		}
		ctx.MiddlewareChains[middlewareCount_global+middlewareCount_router] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func (router *JokerRouter) MapPost(pattern string, handle func(request *http.Request, body []byte, params url.Values, setHeaders func(key, value string)) (status int, response interface{})) {
	pattern = router.prefix + pattern
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount_global := len(router.engine.middlewares)
		middlewareCount_router := len(router.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount_global+middlewareCount_router+1),
			index:            -1,
			maxIndex:         middlewareCount_global + middlewareCount_router + 1,
		}
		copy(ctx.MiddlewareChains, router.engine.middlewares)
		copy(ctx.MiddlewareChains[middlewareCount_global:], router.middlewares)
		finalHandler := func(ctx *JokerContex) {
			if r.Method != http.MethodPost {
				w.WriteHeader(405)
				return
			}
			body := make([]byte, r.ContentLength)
			_, err := r.Body.Read(body)
			if err != nil {
				log.Println("[Error]:Handle in " + pattern + " >>> " + err.Error())
				return
			}
			defer r.Body.Close()
			params := r.URL.Query()
			status, response := handle(r, body, params, func(key, value string) {
				w.Header().Set(key, value)
			})
			if response == nil {
				w.WriteHeader(status)
				return
			}
			jsonresult, err := json.Marshal(response)
			if err != nil {
				log.Println("[Error]:Handle in " + pattern + " >>> " + err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Server", "JokerHttp")
			w.WriteHeader(status)
			w.Write(jsonresult)
		}
		ctx.MiddlewareChains[middlewareCount_global+middlewareCount_router] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func (router *JokerRouter) MapRedirect(pattern string, target string) {
	pattern = router.prefix + pattern
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount_global := len(router.engine.middlewares)
		middlewareCount_router := len(router.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount_global+middlewareCount_router+1),
			index:            -1,
			maxIndex:         middlewareCount_global + middlewareCount_router + 1,
		}
		copy(ctx.MiddlewareChains, router.engine.middlewares)
		copy(ctx.MiddlewareChains[middlewareCount_global:], router.middlewares)
		finalHandler := func(ctx *JokerContex) {
			http.Redirect(w, r, target, http.StatusFound)
		}
		ctx.MiddlewareChains[middlewareCount_global+middlewareCount_router] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func (router *JokerRouter) MapReverseProxy(pattern string, target string) {
	pattern = router.prefix + pattern
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount_global := len(router.engine.middlewares)
		middlewareCount_router := len(router.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount_global+middlewareCount_router+1),
			index:            -1,
			maxIndex:         middlewareCount_global + middlewareCount_router + 1,
		}
		copy(ctx.MiddlewareChains, router.engine.middlewares)
		copy(ctx.MiddlewareChains[middlewareCount_global:], router.middlewares)
		finalHandler := func(ctx *JokerContex) {
			proxy, err := newProxy(target)
			if err != nil {
				ctx.ResponseWriter.WriteHeader(500)
				log.Println("[Error]:Handle in " + pattern + " >>> " + err.Error())
				return
			}
			proxy.ServeHTTP(ctx.ResponseWriter, ctx.Request)
		}
		ctx.MiddlewareChains[middlewareCount_global+middlewareCount_router] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}
