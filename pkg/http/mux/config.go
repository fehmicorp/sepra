package mux

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Mux struct {
	routes map[string]map[string]HandlerFunc
}
