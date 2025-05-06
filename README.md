# JokerHTTP 🃏 - 轻量级 Go HTTP 引擎

![Go 版本](https://img.shields.io/badge/Go-1.16+-blue.svg)
[![许可证](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

<p aligen="center">中文简体 | <a href="README_en.md">English</a></p>

JokerHTTP 是一个轻量级、灵活的 Go HTTP 引擎，让 Web 开发变得简单有趣！🎉

## 功能特性 ✨

- 🚀 支持中间件的简易路由
- 🔥 内置缓存系统
- 📦 静态文件服务
- 🔄 反向代理功能
- ⏱️ 自动缓存过期
- 🛡️ 类型安全处理器
- 🧩 可扩展的中间件架构

## 安装 📦

```bash
go get github.com/jeanhua/jokerhttp
```

## 快速开始 🚀

```go
package main

import (
	"github.com/jeanhua/jokerhttp"
	"net/http"
)

func main() {
	// 创建新的 JokerHTTP 引擎
	engine := jokerhttp.NewEngine()
	engine.Init()
	engine.SetPort(8080)

	// 添加简单的 GET 路由
	engine.MapGet("/hello", func(r *http.Request, params url.Values, setHeaders func(key, value string)) (int, interface{}) {
		return http.StatusOK, map[string]string{"message": "Hello, JokerHTTP! 🎭"}
	})

	// 提供静态文件服务
	engine.UseStaticFiles("./public", "/static")

	// 启动服务器
	engine.Run()
}
```

## API 参考 📚

### 引擎方法

- `Init()` - 使用默认设置初始化引擎
- `SetPort(port int)` - 设置服务器端口
- `Use(middleware Middleware)` - 添加中间件到链中
- `Run()` - 启动服务器

### 路由方法

- `Map(pattern string, handler)` - 通用路由处理器
- `MapGet(pattern string, handler)` - GET 路由处理器
- `MapPost(pattern string, handler)` - POST 路由处理器
- `MapRedirect(pattern string, target string)` - 重定向路由
- `MapReverseProxy(pattern string, target string)` - 反向代理路由

### 缓存方法

- `Set(key string, value interface{}, expiresAt int64)` - 设置缓存值
- `TryGet(key string)` - 获取缓存值
- `Remove(key string)` - 移除缓存项
- `Clear()` - 清除所有缓存
- `AbsoluteTimeFromNow(duration time.Duration)` - 过期时间辅助方法

## 中间件示例 🧩

```go
func LoggerMiddleware(ctx *engine.JokerContex) {
    start := time.Now()
    ctx.Next()
    duration := time.Since(start)
    log.Printf("%s %s - %v", ctx.Request.Method, ctx.Request.URL.Path, duration)
}

// 使用方式:
engine.Use(LoggerMiddleware)
```

## 缓存示例 💾

```go
// 设置5分钟后过期的缓存
expiration := engine.Cache.AbsoluteTimeFromNow(5 * time.Minute)
engine.Cache.Set("user:123", userData, expiration)

// 从缓存获取
if value, ok := engine.Cache.TryGet("user:123"); ok {
    // 使用缓存值
}
```

## 贡献指南 🤝

欢迎贡献！请提交 issue 或 pull request。

## 许可证 📜

MIT 许可证 - 详见 [LICENSE](./LICENSE) 文件。

---

©jeanhua 始于 2025