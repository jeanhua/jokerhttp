# JokerHTTP 🃏 - A Lightweight Go HTTP Engine

![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

JokerHTTP is a lightweight and flexible HTTP engine for Go that makes web development simple and fun! 🎉

## Features ✨

- 🚀 Easy routing with middleware support
- 🔥 Built-in caching system
- 📦 Static file serving
- 🔄 Reverse proxy capabilities
- ⏱️ Automatic cache expiration
- 🛡️ Type-safe handlers
- 🧩 Extensible middleware architecture

## Installation 📦

```bash
go get github.com/jeanhua/jokerhttp
```

## Quick Start 🚀

```go
package main

import (
	"github.com/jeanhua/jokerhttp"
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
- `Use(middleware Middleware)` - Add middleware to the chain
- `Run()` - Start the server

### Routing Methods

- `Map(pattern string, handler)` - Generic route handler
- `MapGet(pattern string, handler)` - GET route handler
- `MapPost(pattern string, handler)` - POST route handler
- `MapRedirect(pattern string, target string)` - Redirect route
- `MapReverseProxy(pattern string, target string)` - Reverse proxy route

### Cache Methods

- `Set(key string, value interface{}, expiresAt int64)` - Set cache value
- `TryGet(key string)` - Get cache value
- `Remove(key string)` - Remove cache item
- `Clear()` - Clear all cache
- `AbsoluteTimeFromNow(duration time.Duration)` - Helper for expiration time

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
// Set cache that expires in 5 minutes
expiration := engine.Cache.AbsoluteTimeFromNow(5 * time.Minute)
engine.Cache.Set("user:123", userData, expiration)

// Get from cache
if value, ok := engine.Cache.TryGet("user:123"); ok {
    // Use cached value
}
```

## Contributing 🤝

Contributions are welcome! Please open an issue or submit a pull request.

## License 📜

MIT License - See [LICENSE](./LICENSE) for details.

---

©jeanhua since 2025