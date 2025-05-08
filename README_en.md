# JokerHTTP 🃏 - Lightweight Go HTTP Engine

![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

<p align="center">English | <a href="README.md">中文简体</a></p>

JokerHTTP is a lightweight, flexible Go HTTP engine that makes web development simple and fun! 🎉

## Features ✨

- 🚀 Simple routing with middleware support
- 🔥 Built-in caching system
- 📦 Static file serving
- 🔄 Reverse proxy functionality
- ⏱️ Automatic cache expiration
- 🛡️ Type-safe handlers
- 🧩 Extensible middleware architecture

## Installation 📦

```bash
go get -u github.com/jeanhua/jokerhttp
```

## Quick Start 🚀

```go
package main

import (
    "github.com/jeanhua/jokerhttp/engine"
	"net/http"
)

func main() {
	// Create a new JokerHTTP engine
	engine := jokerhttp.NewEngine()
	engine.Init()
	engine.SetPort(8080)

	// Add a simple GET route
	engine.MapGet("/hello", func(r *http.Request, params url.Values, setHeaders func(key, value string)) (int, interface{}) {
		return http.StatusOK, map[string]string{"message": "Hello, JokerHTTP! 🎭"}
	})

	// Serve static files
	engine.UseStaticFiles("./public", "/static")

	// Start the server
	engine.Run()
}
```

## API Reference 📚

### Engine Methods

- `Init()` - Initialize the engine with default settings
- `SetPort(port int)` - Set the server port
- `Use(middleware Middleware)` - Add a middleware to the chain
- `Run()` - Start the server

### Router Methods

- `Map(pattern string, handler)` - Generic route handler
- `MapGet(pattern string, handler)` - GET route handler
- `MapPost(pattern string, handler)` - POST route handler
- `MapRedirect(pattern string, target string)` - Redirect route
- `MapReverseProxy(pattern string, target string)` - Reverse proxy route

### Cache Methods

- `Set(key string, value interface{}, expiresAt int64)` - Set a cache value
- `TryGet(key string)` - Get a cached value
- `Remove(key string)` - Remove a cache item
- `Clear()` - Clear all cached items
- `AbsoluteTimeFromNow(duration time.Duration)` - Helper for calculating expiration time

## Middleware Example 🧩

```go
func LoggerMiddleware(ctx *engine.JokerContex) {
    start := time.Now()
    ctx.Next()
    duration := time.Since(start)
    log.Printf("%s %s - %v", ctx.Request.Method, ctx.Request.URL.Path, duration)
}

// Usage:
engine.Use(LoggerMiddleware)
```

## Cache Example 💾

```go
// Set a cache entry that expires in 5 minutes
expiration := utils.AbsoluteTimeFromNow(5 * time.Minute)
engine.Cache.Set("user:123", userData, expiration)

// Retrieve from cache
if value, ok := engine.Cache.TryGet("user:123"); ok {
    // Use cached value
}
```

## Routing Example 🌐

Here's a complete example showing routing with route groups and middleware:

```go
package main

import (
    "github.com/jeanhua/jokerhttp/engine"
    "net/http"
    "net/url"
)

func main() {
    // Initialize the engine
    joker := jokerhttp.NewEngine()
    joker.Init()
    joker.SetPort(1314)

    // Create a router
    router := joker.NewRouter()

    // Root route group
    root := router.Group("/")
    root.Use(func(ctx *engine.JokerContex) {
        ctx.ResponseWriter.Header().Add("middleware", "root")
        ctx.Next()
    })

    // Subgroup /api1
    api1 := root.Group("/api1")
    api1.Use(func(ctx *engine.JokerContex) {
        ctx.ResponseWriter.Header().Add("middleware", "api1")
        ctx.Next()
    })
    api1.Map("/test", func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{}) {
        return 200, "api1 test"
    })

    // Subgroup /api2
    api2 := root.Group("/api2")
    api2.Use(func(ctx *engine.JokerContex) {
        ctx.ResponseWriter.Header().Add("middleware", "api2")
        ctx.Next()
    })
    api2.Map("/test", func(request *http.Request, params url.Values, setHeaders func(key, value string)) (status int, response interface{}) {
        return 200, "api2 test"
    })

    // Start the server
    joker.Run()
}
```

## Contributing 🤝

Contributions are welcome! Please feel free to open an issue or submit a pull request.

## License 📜

MIT License - See [LICENSE](./LICENSE) file for details.

---

©jeanhua Since 2025