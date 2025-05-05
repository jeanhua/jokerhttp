# 🃏 JokerHTTP - 轻量级 Go Web 框架

![Go](https://img.shields.io/badge/Go-1.18%2B-blue)
![license](https://img.shields.io/badge/License-MIT-green)

<p aligen="center">中文简体 | <a href="README_en.md">English</a></p>

JokerHTTP 是一个轻量灵活的 Go Web 框架，旨在让 Web 开发变得简单愉快。🚀

## ✨ 特性

- 🛠️ **中间件支持**：轻松为路由添加中间件
- ⚡ **内置缓存**：简单的内存缓存系统
- 📂 **静态文件**：便捷地提供静态文件服务
- 🔄 **反向代理**：内置反向代理功能
- 🔍 **路由处理**：简单的 GET/POST 路由映射
- ⏱️ **自动缓存清理**：后台 goroutine 清理过期项
- 🔗 **URL 重定向**：轻松实现路由重定向

## 🚀 快速开始

### 安装

```bash
go get github.com/jeanhua/jokerhttp
```

### 基础用法

```go
package main

import (
	"github.com/jeanhua/jokerhttp/engine"
)

func main() {
	// 创建新引擎
	app := engine.NewEngine()
	
	// 使用默认设置初始化
	app.Init()
	
	// 设置自定义端口（默认：9099）
	app.SetPort(8080)
	
	// 添加简单 GET 路由
	app.MapGet("/hello", func(req *http.Request, params url.Values) (int, interface{}) {
		return 200, map[string]string{"message": "你好，JokerHTTP! 👋"}
	})
	
	// 启动服务器
	app.Run()
}
```

## 📚 文档

### 🛠️ 中间件

```go
// 自定义中间件
func LoggerMiddleware(ctx *engine.JokerContex) {
	log.Println("收到请求:", ctx.Request.URL.Path)
	ctx.Next()
}

// 注册中间件
app.Use(LoggerMiddleware)
```

### 💾 缓存使用

```go
// 设置缓存
expireTime := app.Cache.AbsoluteTimeFromNow(5 * time.Minute)
app.Cache.Set("my_key", "my_value", expireTime)

// 获取缓存
if value, ok := app.Cache.TryGet("my_key"); ok {
    fmt.Println("缓存值:", value)
}
```

### 📂 静态文件

```go
// 从 ./public 目录提供静态文件服务，映射到 /static 路径
app.UseStaticFiles("./public", "/static")
```

### 🔄 反向代理

```go
// 将所有 /api 请求代理到另一台服务器
app.MapReverseProxy("/api", "http://api.example.com")
```

### 完整示例：

[examples/basic/main.go](examples/basic/main.go)

## 🤝 贡献指南

欢迎贡献！请提交 issue 或 pull request。

## 📜 许可证

MIT 许可证 - 详见 LICENSE 文件。

------

©Since 2025 jeanhua