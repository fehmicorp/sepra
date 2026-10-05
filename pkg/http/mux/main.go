package mux

import (
	"net/http"
	"strings"
)

func New() *Mux {
	return &Mux{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

func (m *Mux) Handle(method, path string, handler HandlerFunc) {
	method = strings.ToUpper(method)
	if m.routes[method] == nil {
		m.routes[method] = make(map[string]HandlerFunc)
	}
	m.routes[method][path] = handler
}

func (m *Mux) HandleFunc(method, path string, handler HandlerFunc) {
	m.Handle(method, path, handler)
}

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
	handler(w, r)
}
