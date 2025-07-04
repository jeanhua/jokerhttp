package test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jeanhua/jokerhttp/engine"
)

func TestRouter(t *testing.T) {
	joker := engine.NewEngine()
	joker.Init()
	joker.SetPort(1314)
	router := joker.NewRouter()
	root := router.Group("/")
	root.Use(func(ctx *engine.JokerContex) {
		ctx.ResponseWriter.Header().Add("middleware", "root")
		ctx.Next()
	})
	api1 := root.Group("/api1")
	api2 := root.Group("/api2")
	api1.Use(func(ctx *engine.JokerContex) {
		ctx.ResponseWriter.Header().Add("middleware", "api1")
		ctx.Next()
	})
	api2.Use(func(ctx *engine.JokerContex) {
		ctx.ResponseWriter.Header().Add("middleware", "api2")
		ctx.Next()
	})
	api1.Map("/test", func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{}) {
		return 200, "api1 test"
	})
	api2.Map("/test", func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{}) {
		return 200, "api2 test"
	})
	joker.Run()
}

func TestRouterReirect(t *testing.T) {
	joker := engine.NewEngine()
	joker.Init()
	joker.SetPort(1314)
	router := joker.NewRouter()
	root := router.Group("/api")
	root.MapRedirect("/redirect", "https://www.baidu.com")
	joker.Run()
}

func TestRouterProxy(t *testing.T) {
	joker := engine.NewEngine()
	joker.Init()
	joker.SetPort(1314)
	router := joker.NewRouter()
	root := router.Group("/api")
	root.MapReverseProxy("/", "https://imarket.jeanhua.cn/")
	joker.Run()
}

func TestRouterParam(t *testing.T) {
	joker := engine.NewEngine()
	joker.Init()
	joker.SetPort(1314)
	router := joker.NewRouter()
	root := router.Group("/api")
	root.Map("/test/", func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{}) {
		path := strings.Split(request.URL.Path, "/")
		if len(path) < 4 {
			return 400, "error"
		}
		return 200, strings.Split(request.URL.Path, "/")[3]
	})
	joker.Run()
}
