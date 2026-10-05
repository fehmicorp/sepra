package mux

import (
	"net/http"
	"strings"
)

// New initializes and returns a new custom Mux instance.
func New() *Mux {
	return &Mux{
		routes:      make(map[string]map[string]HandlerFunc),
		middlewares: []MiddlewareFunc{},
	}
}

// Use adds optional middleware to the Mux stack safely.
func (m *Mux) Use(middleware MiddlewareFunc) {
	if middleware != nil {
		m.middlewares = append(m.middlewares, middleware)
	}
}

// Handle registers a new handler for the given method and path pattern.
func (m *Mux) Handle(method, path string, handler HandlerFunc) {
	method = strings.ToUpper(method)
	if m.routes[method] == nil {
		m.routes[method] = make(map[string]HandlerFunc)
	}
	m.routes[method][path] = handler
}

// HTTP Method Shortcut Functions

func (m *Mux) GET(path string, handler HandlerFunc)     { m.Handle("GET", path, handler) }
func (m *Mux) POST(path string, handler HandlerFunc)    { m.Handle("POST", path, handler) }
func (m *Mux) PUT(path string, handler HandlerFunc)     { m.Handle("PUT", path, handler) }
func (m *Mux) DELETE(path string, handler HandlerFunc)  { m.Handle("DELETE", path, handler) }
func (m *Mux) PATCH(path string, handler HandlerFunc)   { m.Handle("PATCH", path, handler) }
func (m *Mux) OPTIONS(path string, handler HandlerFunc) { m.Handle("OPTIONS", path, handler) }
func (m *Mux) HEAD(path string, handler HandlerFunc)    { m.Handle("HEAD", path, handler) }
func (m *Mux) TRACE(path string, handler HandlerFunc)   { m.Handle("TRACE", path, handler) }
func (m *Mux) CONNECT(path string, handler HandlerFunc) { m.Handle("CONNECT", path, handler) }

// Routes registers a slice of routes in batch.
func (m *Mux) Routes(routes []Route) error {
	for _, route := range routes {
		m.Handle(route.Method, route.Path, route.Handler)
	}
	return nil
}

// Health is a built-in default health check handler.
func (m *Mux) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ServeHTTP satisfies the standard http.Handler interface and executes middleware chain.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodRoutes, exists := m.routes[r.Method]
	if !exists {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	handler, exists := methodRoutes[r.URL.Path]
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Wrap the core handler with all registered middlewares
	finalHandler := handler
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		if m.middlewares[i] != nil {
			finalHandler = m.middlewares[i](finalHandler)
		}
	}

	finalHandler(w, r)
}

// StartServer initializes the mux, registers built-ins, applies optional middleware, and starts listening.
func StartServer(port string, middleware MiddlewareFunc, routes []Route) error {
	r := New()
	r.GET("/health", r.Health)
	r.Use(middleware)
	r.Routes(routes)
	return http.ListenAndServe(":"+port, r)
}
