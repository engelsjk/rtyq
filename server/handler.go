package server

import (
	"fmt"
	"net/http"

	"github.com/engelsjk/rtyq/conf"
	"github.com/engelsjk/rtyq/data"
	"github.com/go-chi/chi"
)

const (
	routeVarWildcard = "*"
	routeVarLayer    = "layer"
	routeVarSubLayer = "sublayer"
	routeVarPoint    = "point"
	routeVarBBox     = "bbox"
	routeVarID       = "id"
	routeVarTileX    = "x"
	routeVarTileY    = "y"
	routeVarTileZ    = "z"
)

var (
	ErrNotFound error = fmt.Errorf("not found")
)

func initRouter() *chi.Mux {
	router := chi.NewRouter()
	return router
}

func addRoutes(router *chi.Mux) {
	addRoute(router, "/", handleRoot)

	addRoute(router, "/*", handleDefault)

	addRoute(router, "/{layer}", handleLayer)
	addRoute(router, "/{layer}/{sublayer}", handleLayer)
	addRoute(router, "/{layer}/{sublayer}/*", handleLayer)

	addRoute(router, "/{layer}/point", handlePoint)
	addRoute(router, "/{layer}/{sublayer}/point", handlePoint)
	addRoute(router, "/{layer}/point/{point}", handlePoint)
	addRoute(router, "/{layer}/{sublayer}/point/{point}", handlePoint)

	addRoute(router, "/{layer}/bbox", handleBBox)
	addRoute(router, "/{layer}/{sublayer}/bbox", handleBBox)
	addRoute(router, "/{layer}/bbox/{bbox}", handleBBox)
	addRoute(router, "/{layer}/{sublayer}/bbox/{bbox}", handleBBox)

	addRoute(router, "/{layer}/tile", handleTile)
	addRoute(router, "/{layer}/{sublayer}/tile", handleTile)
	addRoute(router, "/{layer}/tile/{z}", handleTile)
	addRoute(router, "/{layer}/{sublayer}/tile/{z}", handleTile)
	addRoute(router, "/{layer}/tile/{z}/{x}", handleTile)
	addRoute(router, "/{layer}/{sublayer}/tile/{z}/{x}", handleTile)
	addRoute(router, "/{layer}/tile/{z}/{x}/{y}", handleTile)
	addRoute(router, "/{layer}/{sublayer}/tile/{z}/{x}/{y}", handleTile)

	addRoute(router, "/{layer}/id", handleID)
	addRoute(router, "/{layer}/{sublayer}/id", handleID)
	addRoute(router, "/{layer}/id/{id}", handleID)
	addRoute(router, "/{layer}/{sublayer}/id/{id}", handleID)

	addRoute(router, "/config", handleConfig)
}

func addRoute(router *chi.Mux, path string, handler func(http.ResponseWriter, *http.Request) *serverError) {
	router.Handle(path, serverHandler(handler))
}

////////////////////////////////////////////////////////////////////////

func handleRoot(w http.ResponseWriter, r *http.Request) *serverError {
	type Home struct {
		API    string `json:"api"`
		Config string `json:"config"`
	}

	api := fmt.Sprintf("%s %s", conf.AppConfig.Name, conf.AppConfig.Version)

	home := Home{
		API:    api,
		Config: "/config",
	}
	return writeJSON(w, ContentTypeJSON, home)
}

func handleDefault(w http.ResponseWriter, r *http.Request) *serverError {
	return serverErrorNotFound(ErrNotFound, ErrNotFound.Error())
}

func handleLayer(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	sublayer := getRequestVar(routeVarSubLayer, r)
	wildcard := getRequestVar(routeVarWildcard, r)

	if layer == "" {
		return errorQueryToServer(data.ErrQueryMissingLayer)
	}

	if sublayer != "" {
		layer = fmt.Sprintf("%s/%s", layer, sublayer)
	}

	if !data.QueryHandler.HasLayer(layer) {
		return errorQueryToServer(data.ErrQueryInvalidLayer)
	}

	if wildcard == "" {
		return errorQueryToServer(data.ErrQueryMissingQuery)
	}

	return errorQueryToServer(data.ErrQueryInvalidQuery)
}

func handlePoint(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	sublayer := getRequestVar(routeVarSubLayer, r)
	point := getRequestVar(routeVarPoint, r)

	if sublayer != "" {
		layer = fmt.Sprintf("%s/%s", layer, sublayer)
	}

	features, err := data.QueryHandler.Point(layer, point)
	if err != nil {
		return errorQueryToServer(err)
	}

	return writeJSON(w, ContentTypeJSON, features)
}

func handleBBox(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	sublayer := getRequestVar(routeVarSubLayer, r)
	bbox := getRequestVar(routeVarBBox, r)

	if sublayer != "" {
		layer = fmt.Sprintf("%s/%s", layer, sublayer)
	}

	features, err := data.QueryHandler.BBox(layer, bbox)
	if err != nil {
		return errorQueryToServer(err)
	}

	return writeJSON(w, ContentTypeJSON, features)
}

func handleTile(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	sublayer := getRequestVar(routeVarSubLayer, r)
	tileX := getRequestVar(routeVarTileX, r)
	tileY := getRequestVar(routeVarTileY, r)
	tileZ := getRequestVar(routeVarTileZ, r)

	if sublayer != "" {
		layer = fmt.Sprintf("%s/%s", layer, sublayer)
	}

	features, err := data.QueryHandler.Tile(layer, tileX, tileY, tileZ)
	if err != nil {
		return errorQueryToServer(err)
	}

	return writeJSON(w, ContentTypeJSON, features)
}

func handleID(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	sublayer := getRequestVar(routeVarSubLayer, r)
	id := getRequestVar(routeVarID, r)

	if sublayer != "" {
		layer = fmt.Sprintf("%s/%s", layer, sublayer)
	}

	features, err := data.QueryHandler.ID(layer, id)
	if err != nil {
		return errorQueryToServer(err)
	}

	return writeJSON(w, ContentTypeJSON, features)
}

func handleConfig(w http.ResponseWriter, r *http.Request) *serverError {

	type Config struct {
		Layers  []string `json:"layers"`
		Queries []string `json:"queries"`
	}

	layers := data.QueryHandler.Layers()

	queries := []string{
		"/{layer}/point/{point}",
		"/{layer}/tile/{z}/{x}/{y}",
		"/{layer}/id/{id}",
	}

	config := Config{Layers: layers, Queries: queries}

	return writeJSON(w, ContentTypeJSON, config)
}

/////////////////////////////////////////////////////

func errorQueryToServer(err error) *serverError {
	switch err {
	case data.ErrQueryMissingLayer:
		return serverErrorNotFound(err, err.Error())
	case data.ErrQueryInvalidLayer:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryMissingQuery:
		return serverErrorNotFound(err, err.Error())
	case data.ErrQueryInvalidQuery:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryMissingID:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryInvalidID:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryMissingPoint:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryInvalidPoint:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryMissingTile:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryInvalidTile:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryMissingBBox:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryInvalidBBox:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryExceededTileZoomLimit:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryRequest:
		return serverErrorInternal(err, err.Error())
	default:
		err := fmt.Errorf("unknown error")
		return serverErrorInternal(err, err.Error())
	}
}
