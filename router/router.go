package router

import (
	"github.com/gorilla/mux"
)

// NewRouter builds and returns a new router from routes
func NewRouter() *mux.Router {
	// When StrictSlash == true, if the route path is "/path/", accessing "/path" will perform a redirect to the former and vice versa.
	router := mux.NewRouter().StrictSlash(true)
	router.Use(Logger)
	router.Use(BasicAuth)

	sub := router.PathPrefix("/api/v1").Subrouter()

	// Append routes for all the objects
	routes := routes0
	routes = append(routes, routes1...)
	routes = append(routes, routes2...)
	routes = append(routes, routes3...)
	routes = append(routes, routes4...)
	routes = append(routes, routes5...)
	routes = append(routes, routes6...)
	routes = append(routes, routes7...)

	for _, route := range routes {
		sub.
			HandleFunc(route.Pattern, route.HandlerFunc).
			Name(route.Name).
			Methods(route.Method)
	}

	return router
}
