# server

Starts an HTTP server with sensible timeouts and graceful shutdown on `SIGINT`/`SIGTERM`.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users", getUsers)
mux.HandleFunc("POST /users", createUser)

// Blocks until the process receives a shutdown signal,
// then drains connections with a 15-second grace period.
server.ServeWithGracefulShutdown(":8080", mux)
```
