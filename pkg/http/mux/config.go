package mux

import "net/http"

// HandlerFunc defines the signature for request handlers.
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// MiddlewareFunc defines the signature for middleware wrapping HandlerFunc.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// Route represents a single route definition for batch registration.
type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

// Mux is a custom HTTP request multiplexer.
type Mux struct {
	routes      map[string]map[string]HandlerFunc
	middlewares []MiddlewareFunc
}
