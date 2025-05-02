# 🃏 JokerHttp —— 轻量级 Go HTTP 框架

## 🔍 概述

JokerHttp 是一个轻量、简洁的 Go 语言 HTTP 框架，提供如下核心功能：

- 灵活的路由映射
- 支持中间件机制
- 静态文件服务
- 上下文管理（Context）

适用于快速搭建 Web 服务或 API 接口。

---

## 🚀 快速开始

### 安装

使用 `go get` 安装 JokerHttp：

```bash
go get "github.com/jeanhua/jokerhttp"
```

### 基础示例

```go
package main

import (
	"github.com/jeanhua/jokerhttp"
	"github.com/jeanhua/jokerhttp/engine"
)

func main() {
	// 创建引擎实例
	engine := jokerhttp.NewEngine()
	engine.Init()

	// 设置端口（可选，默认为 9099）
	engine.SetPort(8080)

	// 添加路由
	engine.MapGet("/hello", helloHandler)
	engine.MapPost("/echo", echoHandler)

	// 启动服务
	engine.Run()
}

func helloHandler(request *http.Request, params url.Values) (int, interface{}) {
	return 200, map[string]string{"message": "Hello, World!"}
}

func echoHandler(request *http.Request, body []byte, params url.Values) (int, interface{}) {
	return 200, map[string]string{"echo": string(body)}
}
```

---

## 🛠 功能详解

### 1. 初始化引擎

```go
engine := jokerhttp.NewEngine()
engine.Init() // 默认端口为 9099
```

### 2. 设置端口

```go
engine.SetPort(8080) // 自定义监听端口
```

### 3. 静态文件服务

```go
engine.UseStaticFiles("./public") // 设置静态资源根目录
```

### 4. 路由映射

#### GET 请求

```go
engine.MapGet("/path", func(r *http.Request, p url.Values) (int, interface{}) {
    return http.StatusOK, map[string]string{"key": "value"}
})
```

#### POST 请求

```go
engine.MapPost("/path", func(r *http.Request, b []byte, p url.Values) (int, interface{}) {
    return http.StatusOK, map[string]string{"received": string(b)}
})
```

### 5. 中间件支持

#### 添加中间件

```go
engine.Use(func(ctx *engine.Contex) {
    fmt.Println("Before handler")
    ctx.Next()
    fmt.Println("After handler")
})
```

#### 控制中间件流程

```go
engine.Use(func(ctx *engine.Contex) {
    if !checkAuth(ctx.Request) {
        ctx.AbortWithStatusJSON(401, map[string]string{"error": "Unauthorized"})
        return
    }
    ctx.Next()
})
```

### 6. Context 方法一览

| 方法 | 描述 |
|------|------|
| `ctx.Next()` | 执行下一个中间件或处理函数 |
| `ctx.Abort()` | 终止中间件链执行 |
| `ctx.AbortWithStatus(statusCode int)` | 终止并返回指定状态码 |
| `ctx.AbortWithStatusJSON(statusCode int, jsonObj interface{})` | 终止并返回 JSON 格式响应 |
| `ctx.Use(middleware Middleware)` | 动态添加中间件 |

---

## 💡 示例代码：完整用法演示

```go
package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jeanhua/jokerhttp"
	"github.com/jeanhua/jokerhttp/engine"
)

func main() {
	joker := jokerhttp.NewEngine()
	joker.Init()
	joker.SetPort(1314)

	// 日志中间件
	joker.Use(func(ctx *engine.Contex) {
		println("Request received:", ctx.Request.URL.Path)
		ctx.Next()
		println("Response sent")
	})

	// 认证中间件
	joker.Use(func(ctx *engine.Contex) {
		if ctx.Request.Header.Get("Authorization") != "secret" {
			ctx.AbortWithStatusJSON(401, map[string]string{"error": "Unauthorized"})
			return
		}
		ctx.Next()
	})

	// 路由设置
	joker.MapGet("/hello", func(r *http.Request, p url.Values) (int, interface{}) {
		name := p.Get("name")
		if name == "" {
			name = "World"
		}
		return 200, map[string]string{"message": "Hello, " + name + "!"}
	})

	joker.MapPost("/echo", func(r *http.Request, b []byte, p url.Values) (int, interface{}) {
		return 200, struct {
			Original string `json:"original"`
		}{
			Original: string(b),
		}
	})

	// 静态文件服务
	joker.UseStaticFiles("./static")

	// 启动服务
	fmt.Println("服务启动...")
	joker.Run()
}
```

> ✅ **提示**：如果不使用map，你可以返回任意结构体作为响应数据，框架会自动进行 JSON 序列化。

---

## 🧪 测试接口

### GET 请求测试

```bash
curl "http://localhost:8080/hello?name=Joker"
```

### POST 请求测试

```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: secret" \
     -d '{"key":"value"}' \
     http://localhost:8080/echo
```

---

## ⚠️ 注意事项

1. **默认端口**为 `9099`，可通过 `SetPort()` 修改。
2. **中间件顺序**很重要，按添加顺序依次执行。
3. 使用 `Abort()` 或其变种方法可提前终止请求处理链。
4. 静态文件服务默认映射到 `/` 路径。

---

如需进一步了解，请查看 [GitHub 项目地址](https://github.com/jeanhua/jokerhttp) 获取最新文档和更新。欢迎贡献代码或提出建议！

---

