package test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/jeanhua/jokerhttp/engine"
)

func TestServer(t *testing.T) {
	server := &engine.JokerEngine{}
	server.Init()
	// test get localhost:9909/json
	server.MapGet("/json", backjson)
	// test get localhost:9909/int
	server.MapGet("/int", backint)
	// test get localhost:9909/string
	server.MapGet("/string", backString)

	server.Run()
}

func backint(request *http.Request, params url.Values, setHeader func(key, value string)) (status int, response interface{}) {
	return 200, 114514
}

func backjson(request *http.Request, params url.Values, setHeader func(key, value string)) (status int, response interface{}) {
	return 200, struct {
		Message string
	}{
		Message: "success",
	}
}

func backString(request *http.Request, params url.Values, setHeader func(key, value string)) (status int, response interface{}) {
	return 200, "success"
}
