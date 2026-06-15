package api

import "net/http"

// RouteRegistrar abstracts the underlying chi.Router so that feature handlers
// don't need to depend on chi directly.
type RouteRegistrar interface {
	Get(pattern string, handlerFn http.HandlerFunc)
	Post(pattern string, handlerFn http.HandlerFunc)
	Put(pattern string, handlerFn http.HandlerFunc)
	Delete(pattern string, handlerFn http.HandlerFunc)
}

// Handler is implemented by each feature package to self-register its routes.
type Handler interface {
	RegisterRoutes(r RouteRegistrar)
}
