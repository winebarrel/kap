package kap

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"slices"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

type AuthHandler struct {
	Options *Options
	Proxy   func(http.ResponseWriter, *http.Request)
}

func NewAuthHandler(options *Options) *AuthHandler {
	proxy := httputil.NewSingleHostReverseProxy(options.Backend)

	handler := &AuthHandler{
		Options: options,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			r.Host = options.Backend.Host
			proxy.ServeHTTP(w, r)
		},
	}

	return handler
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get(h.Options.Key)

	if secret == "" {
		secret = r.URL.Query().Get(h.Options.Key)
	}

	if !slices.Contains(h.Options.Secret, secret) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	h.Proxy(w, r)
}
