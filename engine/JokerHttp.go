package engine

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type JokerEngine struct {
	port        int
	middlewares []Middleware
	Cache       *jokerCache
}

func NewEngine() *JokerEngine {
	return &JokerEngine{}
}

func (jokerEngine *JokerEngine) Init() {
	log.SetFlags(log.LUTC)
	if jokerEngine.port == 0 {
		jokerEngine.port = 9099
	}
	// Initialize the cache
	jokerEngine.Cache = &jokerCache{}
	jokerEngine.Cache.init()
}

func (jokerEngine *JokerEngine) SetPort(port int) {
	jokerEngine.port = port
}

func (jokerEngine *JokerEngine) Use(middleware Middleware) {
	jokerEngine.middlewares = append(jokerEngine.middlewares, middleware)
}

func (jokerEngine *JokerEngine) UseStaticFiles(baseRoot string, target string) {
	baseRoot = strings.ReplaceAll(baseRoot, "\\", "/")
	fs := http.FileServer(http.Dir(baseRoot))
	if _, err := os.Stat(baseRoot); err != nil && os.IsNotExist(err) {
		log.Printf("Directory does not exist: %s\n", baseRoot)
	}
	// Handle the static file server
	http.HandleFunc(target, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "JokerHttp")
		w.Header().Set("X-Static-File", "JokerHttp")
		w.Header().Set("Cache-Control", "cache, max-age=3600")
		if strings.HasPrefix(r.URL.Path, target) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, target)
			fs.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

func (jokerEngine *JokerEngine) Map(pattern string, handle func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{})) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(jokerEngine.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount+1),
			index:            -1,
			maxIndex:         middlewareCount + 1,
			aborted:          false,
		}
		copy(ctx.MiddlewareChains, jokerEngine.middlewares)
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

		ctx.MiddlewareChains[middlewareCount] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func (jokerEngine *JokerEngine) MapGet(pattern string, handle func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{})) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(jokerEngine.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount+1),
			index:            -1,
			maxIndex:         middlewareCount + 1,
			aborted:          false,
		}
		copy(ctx.MiddlewareChains, jokerEngine.middlewares)
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

		ctx.MiddlewareChains[middlewareCount] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func (jokerEngine *JokerEngine) MapPost(pattern string, handle func(request *http.Request, body []byte, params url.Values, setHeaders func(key, value string)) (status int, response interface{})) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(jokerEngine.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount+1),
			index:            -1,
			maxIndex:         middlewareCount + 1,
			aborted:          false,
		}
		copy(ctx.MiddlewareChains, jokerEngine.middlewares)
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
			jsonResult, err := json.Marshal(response)
			if err != nil {
				log.Println("[Error]: Failed to marshal response:", err)
				w.WriteHeader(500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Server", "JokerHttp")
			w.WriteHeader(status)
			w.Write(jsonResult)
		}

		ctx.MiddlewareChains[middlewareCount] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func (jokerEngine *JokerEngine) Run() {
	http.ListenAndServe(":"+strconv.Itoa(jokerEngine.port), nil)
}

func (jokerEngine *JokerEngine) RunWithAddr(addr string) {
	http.ListenAndServe(addr, nil)
}

func (jokerEngine *JokerEngine) MapRedirect(pattern string, target string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(jokerEngine.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount+1),
			index:            -1,
			maxIndex:         middlewareCount + 1,
			aborted:          false,
		}
		copy(ctx.MiddlewareChains, jokerEngine.middlewares)
		finalHandler := func(ctx *JokerContex) {
			http.Redirect(w, r, target, http.StatusFound)
		}
		ctx.MiddlewareChains[middlewareCount] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}

func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = modifyResponse()
	return proxy, nil
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "JokerHttp")
		resp.Header.Set("Server", "JokerHttp")
		return nil
	}
}

func (jokerEngine *JokerEngine) MapReverseProxy(pattern string, target string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(jokerEngine.middlewares)
		ctx := &JokerContex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount+1),
			index:            -1,
			maxIndex:         middlewareCount + 1,
			aborted:          false,
		}
		copy(ctx.MiddlewareChains, jokerEngine.middlewares)
		finalHandler := func(ctx *JokerContex) {
			proxy, err := newProxy(target)
			if err != nil {
				ctx.ResponseWriter.WriteHeader(500)
				return
			}
			proxy.ServeHTTP(ctx.ResponseWriter, ctx.Request)
		}
		ctx.MiddlewareChains[middlewareCount] = finalHandler
		if len(ctx.MiddlewareChains) > 0 {
			ctx.Next()
		}
	})
}
