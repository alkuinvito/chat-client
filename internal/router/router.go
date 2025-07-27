package router

import (
	"chat-client/internal/discovery"
	"fmt"
	"net/http"
)

type RouterConfig struct {
	addr   string
	prefix string
}

type Router struct {
	handler *http.ServeMux
	server  *http.Server
}

type IRouter interface {
	Handle() error
	Use(path string, handler http.Handler) *Router
}

func DefaultConfig() *RouterConfig {
	addr := fmt.Sprintf(":%d", discovery.SVC_PORT)
	prefix := "/api"

	return &RouterConfig{addr, prefix}
}

func NewRouter(config *RouterConfig) *Router {
	addr := config.addr
	handler := http.NewServeMux()

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	if config.prefix != "" {
		withPrefix := http.NewServeMux()
		withPrefix.Handle(config.prefix+"/", http.StripPrefix(config.prefix, handler))

		server.Handler = withPrefix
	}

	return &Router{handler, &server}
}

func (r *Router) Handle() error {
	return r.server.ListenAndServe()
}

func (r *Router) Use(path string, handler http.Handler) *Router {
	r.handler.Handle(path+"/", http.StripPrefix(path, handler))
	return r
}
