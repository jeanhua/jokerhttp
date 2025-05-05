package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jeanhua/jokerhttp"
	"github.com/jeanhua/jokerhttp/engine"
)

func main() {
	// 初始化引擎
	joker := jokerhttp.NewEngine()
	joker.Init()
	joker.SetPort(1314)

	// 添加日志中间件
	joker.Use(func(ctx *engine.Contex) {
		println("Request received:", ctx.Request.URL.Path)
		ctx.Next()
		println("Response sent")
	})

	// 添加认证中间件
	joker.Use(func(ctx *engine.Contex) {
		if ctx.Request.Header.Get("Authorization") != "secret" {
			ctx.AbortWithStatusJSON(401, map[string]string{"error": "Unauthorized"})
			ctx.Abort()
		}
		ctx.Next()
	})

	// 添加路由
	joker.MapGet("/hello", func(r *http.Request, p url.Values) (int, interface{}) {
		name := p.Get("name")
		if name == "" {
			name = "World"
		}
		return 200, map[string]string{"message": "Hello, " + name + "!"}
	})

	joker.MapPost("/echo", func(r *http.Request, b []byte, p url.Values) (int, interface{}) {
		return 200, map[string]string{"original": string(b)}
	})

	// 静态文件服务
	joker.UseStaticFiles("./static", "/")

	// 启动服务
	fmt.Println("http://localhost:1314 -> 服务启动...")
	joker.Run()
}
