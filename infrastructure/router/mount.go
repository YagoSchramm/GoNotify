package router

import (
	"github.com/gorilla/mux"
)

func Mount(root *mux.Router, protectedMiddlewares []mux.MiddlewareFunc, modules ...Module) {
	for _, module := range modules {
		moduleRouter := root.PathPrefix(module.Path()).Subrouter()
		protectedRouter := moduleRouter.NewRoute().Subrouter()

		if len(protectedMiddlewares) > 0 {
			protectedRouter.Use(protectedMiddlewares...)
		}
		if len(module.Middlewares()) > 0 {
			protectedRouter.Use(module.Middlewares()...)
		}

		for _, route := range module.Routes() {
			routeRouter := moduleRouter
			if !route.Public {
				routeRouter = protectedRouter
			}

			if len(route.Middlewares) > 0 {
				routeRouter = routeRouter.NewRoute().Subrouter()
				routeRouter.Use(route.Middlewares...)
			}

			routeRouter.Handle(route.Path, route.Handler).Methods(route.HttpMethods...)
		}
	}
}
