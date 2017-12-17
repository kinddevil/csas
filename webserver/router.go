package webserver

import (
	"github.com/gorilla/mux"
)

func NewRouter(routes *Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range *routes {

		// Attach each route, uses a Builder-like pattern to set each route up.
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func NewRouterWithPrefix(routes *Routes, prefix string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range *routes {

		// Attach each route, uses a Builder-like pattern to set each route up.
		router.Methods(route.Method).
			Path(prefix + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
