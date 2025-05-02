package engine

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type JokerEngine struct {
	port        int
	middlewares []Middleware
}

func (jokerEngine *JokerEngine) Init() {
	log.SetFlags(log.LUTC)
	if jokerEngine.port == 0 {
		jokerEngine.port = 9099
	}
}

func (jokerEngine *JokerEngine) SetPort(port int) {
	jokerEngine.port = port
}

func (jokerEngine *JokerEngine) Use(middleware Middleware) {
	jokerEngine.middlewares = append(jokerEngine.middlewares, middleware)
}

func (jokerEngine *JokerEngine) UseStaticFiles(baseRoot string) {
	baseRoot = strings.ReplaceAll(baseRoot, "\\", "/")
	fs := http.FileServer(http.Dir(baseRoot))
	if _, err := os.Stat(baseRoot); err != nil && os.IsNotExist(err) {
		log.Printf("Directory does not exist: %s\n", baseRoot)
	}
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func (jokerEngine *JokerEngine) MapGet(pattern string, handle func(request *http.Request, params url.Values) (status int, response interface{})) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(jokerEngine.middlewares)
		ctx := &Contex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount+1),
			index:            -1,
			maxIndex:         middlewareCount + 1,
		}
		copy(ctx.MiddlewareChains, jokerEngine.middlewares)
		finalHandler := func(ctx *Contex) {
			if r.Method != http.MethodGet {
				w.WriteHeader(405)
				return
			}
			params := r.URL.Query()
			status, response := handle(r, params)
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

func (jokerEngine *JokerEngine) MapPost(pattern string, handle func(request *http.Request, body []byte, params url.Values) (status int, response interface{})) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		middlewareCount := len(jokerEngine.middlewares)
		ctx := &Contex{
			Request:          r,
			ResponseWriter:   w,
			MiddlewareChains: make([]Middleware, middlewareCount+1),
			index:            -1,
			maxIndex:         middlewareCount + 1,
		}
		copy(ctx.MiddlewareChains, jokerEngine.middlewares)
		finalHandler := func(ctx *Contex) {
			if r.Method != http.MethodPost {
				w.WriteHeader(405)
				return
			}
			var body []byte
			var err error

			if r.Header.Get("Content-Type") == "application/json" {
				body, err = io.ReadAll(r.Body)
				if err != nil {
					log.Println("[Error]: Failed to read request body:", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				defer r.Body.Close()
			} else {
				err := r.ParseForm()
				if err != nil {
					log.Println("[Error]: Failed to parse form:", err)
					w.WriteHeader(400)
					return
				}
			}
			params := r.URL.Query()
			status, response := handle(r, body, params)
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
