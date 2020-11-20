package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

const (
	routeVarLayer = "layer"
	routeVarPoint = "point"
	routeVarTile  = "tile"
	routeVarID    = "id"
)

func initRouter() *chi.Mux {

	router := chi.NewRouter()

	addRoute(router, "/", handleRoot)

	addRoute(router, "/{layer}", handleLayer)

	addRoute(router, "/{layer}/point/{point}", handlePoint)

	addRoute(router, "/{layer}/tile/{tile}", handleTile)

	addRoute(router, "/{layer}/id/{id}", handleID)

	return router
}

func addRoute(router *chi.Mux, path string, handler func(http.ResponseWriter, *http.Request) *serverError) {
	router.Handle(path, serverHandler(handler))
}

func handleRoot(w http.ResponseWriter, r *http.Request) *serverError {
	return nil
}

func handleLayer(w http.ResponseWriter, r *http.Request) *serverError {
	return nil
}

func handlePoint(w http.ResponseWriter, r *http.Request) *serverError {

	_ = getRequestVar(routeVarLayer, r)
	_ = getRequestVar(routeVarPoint, r)

	return nil
}

func handleTile(w http.ResponseWriter, r *http.Request) *serverError {

	_ = getRequestVar(routeVarLayer, r)
	_ = getRequestVar(routeVarTile, r)

	return nil
}

func handleID(w http.ResponseWriter, r *http.Request) *serverError {

	_ = getRequestVar(routeVarLayer, r)
	_ = getRequestVar(routeVarTile, r)

	return nil
}
