package test

import (
	"net/http"
	"net/url"
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
	api := root.Group("/api")
	api.Use(func(ctx *engine.JokerContex) {
		ctx.ResponseWriter.Header().Add("middleware", "api")
		ctx.Next()
	})
	api.Map("/test", func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{}) {
		setHeaders("Content-Type", "application/json")
		return 200, map[string]interface{}{
			"code":    0,
			"message": "ok",
			"data": map[string]interface{}{
				"method": request.Method,
				"path":   request.URL.Path,
				"query":  params,
			},
		}
	})
	joker.Run()
}
