package mux

import (
	"net/http"
	"strings"

	cf "github.com/sepra/pkg/http/config"
)

func New() *cf.Mux {
	return &cf.Mux{
		routes: make(map[string]map[string]cf.HandlerFunc),
	}
}

func (m *cf.Mux) Handle(method, path string, handler cf.HandlerFunc) {
	method = strings.ToUpper(method)
	if m.routes[method] == nil {
		m.routes[method] = make(map[string]cf.HandlerFunc)
	}
	m.routes[method][path] = handler
}

func (m *cf.Mux) HandleFunc(method, path string, handler cf.HandlerFunc) {
	m.Handle(method, path, handler)
}

func (m *cf.Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
