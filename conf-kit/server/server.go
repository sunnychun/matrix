package server

import (
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	op, key, err := parsePath(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("%s %q", req.Method, req.URL.Path)

	switch op {
	case "list":
		h.list(w, req, key)
	case "conf":
		h.conf(w, req, key)
	default:
		http.Error(w, fmt.Sprintf("no such operation: %q", op), http.StatusBadRequest)
	}
}

func (h *Handler) list(w http.ResponseWriter, req *http.Request, key string) {
}

func (h *Handler) conf(w http.ResponseWriter, req *http.Request, key string) {
}
