package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jeanhua/jokerhttp"
	"github.com/jeanhua/jokerhttp/engine"
)

func main() {
	// 初始化引擎
	joker := jokerhttp.NewEngine()
	joker.Init()
	joker.SetPort(1314)

	// 添加日志中间件
	joker.Use(func(ctx *engine.JokerContex) {
		println("Request received:", ctx.Request.URL.Path)
		ctx.Next()
		println("Response sent")
	})

	// 添加认证中间件
	joker.Use(func(ctx *engine.JokerContex) {
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

	joker.MapGet("/cache", func(request *http.Request, params url.Values) (status int, response interface{}) {
		// 获取缓存值
		value, found := joker.Cache.TryGet("myKey")
		if found {
			return 200, map[string]interface{}{"value": value}
		}
		// 设置缓存值
		minutes := time.Now().Minute()
		hours := time.Now().Hour()
		seconds := time.Now().Second()
		joker.Cache.Set("myKey", map[string]interface{}{"time": fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)}, joker.Cache.AbsoluteTimeFromNow(time.Second*60)) // 60秒后过期
		return 200, map[string]interface{}{"message": "Cache set!"}
	})

	// 静态文件服务
	joker.UseStaticFiles("./static", "/")

	// 启动服务
	fmt.Println("http://localhost:1314 -> 服务启动...")
	joker.Run()
}
