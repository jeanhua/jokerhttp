# 🃏 JokerHTTP - A Lightweight Go Web Framework

![Go](https://img.shields.io/badge/Go-1.18%2B-blue)
![license](https://img.shields.io/badge/License-MIT-green)

JokerHTTP is a lightweight and flexible web framework for Go, designed to make web development simple and enjoyable. 🚀

## ✨ Features

- 🛠️ **Middleware Support**: Easily add middleware to your routes
- ⚡ **Built-in Cache**: Simple in-memory caching system
- 📂 **Static Files**: Serve static files with ease
- 🔄 **Reverse Proxy**: Built-in reverse proxy capabilities
- 🔍 **Route Handling**: Simple GET/POST route mapping
- ⏱️ **Automatic Cache Cleanup**: Background goroutine cleans expired items
- 🔗 **URL Redirection**: Easy route redirection

## 🚀 Quick Start

### Installation

```bash
go get github.com/jeanhua/jokerhttp
```

### Basic Usage

```go
package main

import (
	"github.com/jeanhua/jokerhttp/engine"
)

func main() {
	// Create new engine
	app := engine.NewEngine()
	
	// Initialize with default settings
	app.Init()
	
	// Set custom port (default: 9099)
	app.SetPort(8080)
	
	// Add a simple GET route
	app.MapGet("/hello", func(req *http.Request, params url.Values) (int, interface{}) {
		return 200, map[string]string{"message": "Hello, JokerHTTP! 👋"}
	})
	
	// Start the server
	app.Run()
}
```

## 📚 Documentation

### 🛠️ Middleware

```go
// Custom middleware
func LoggerMiddleware(ctx *engine.JokerContex) {
	log.Println("Request received:", ctx.Request.URL.Path)
	ctx.Next()
}

// Register middleware
app.Use(LoggerMiddleware)
```

### 💾 Cache Usage

```go
// Set cache
expireTime := app.Cache.AbsoluteTimeFromNow(5 * time.Minute)
app.Cache.Set("my_key", "my_value", expireTime)

// Get cache
if value, ok := app.Cache.TryGet("my_key"); ok {
    fmt.Println("Cached value:", value)
}
```

### 📂 Static Files

```go
// Serve static files from ./public at /static
app.UseStaticFiles("./public", "/static")
```

### 🔄 Reverse Proxy

```go
// Proxy all requests from /api to another server
app.MapReverseProxy("/api", "http://api.example.com")
```

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## 📜 License

MIT License - see LICENSE for details.

------

©Since 2025 jeanhua