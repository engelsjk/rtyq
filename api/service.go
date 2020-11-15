package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/engelsjk/rtyq"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var (
	ErrUnableToWriteMessage error = fmt.Errorf("unable to write message")
)

// Message acts as the home endpoint output for the api service
type Message struct {
	Message   string   `json:"message"`
	Endpoints []string `json:"endpoints"`
}

// Start initializes routes for each layer (data/database/index)
// and starts an api at the specified port
func Start(cfg *rtyq.Config) error {

	err := rtyq.ValidateConfigData(cfg)
	if err != nil {
		return err
	}

	err = rtyq.ValidateConfigDatabase(cfg)
	if err != nil {
		return err
	}

	err = rtyq.ValidateConfigServiceOnly(cfg)
	if err != nil {
		return err
	}

	router := chi.NewRouter()

	// middleware
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.Throttle(cfg.ThrottleLimit))

	router.Use(middleware.Recoverer)

	if cfg.EnableLogs {
		router.Use(middleware.Logger)
	}

	// add routes for each data layer
	numRoutes, err := SetRoutes(router, cfg)
	if err != nil {
		return err
	}

	plural := ""
	if numRoutes > 1 {
		plural = "s"
	}

	fmt.Println("%************%")
	fmt.Printf("starting %d layer service%s on localhost:%d\n", numRoutes, plural, cfg.Port)

	return http.ListenAndServe(
		net.JoinHostPort("", strconv.Itoa(cfg.Port)),
		router,
	)
}

// SetRoutes initializes all of the API service endpoints. 	/
// It iterates over each layer to initialize each db with a spatial index
// and links it to a separate layer endpoint.
func SetRoutes(router *chi.Mux, cfg *rtyq.Config) (int, error) {

	numRoutes := 0
	layerEndpoints := []string{}

	for ii := 0; ii < len(cfg.Layers); ii++ {

		layer := cfg.Layers[ii]

		fn := filepath.Base(layer.Database.Path)

		fmt.Println("%************%")
		fmt.Printf("setting route for layer: %s\n", layer.Name)

		data, err := rtyq.InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
		if err != nil {
			fmt.Printf("warning : layer (%s) : %s : skipping layer\n", layer.Name, err.Error())
			continue
		}

		if !rtyq.FileExists(layer.Database.Path) {
			fmt.Printf("warning : layer (%s) : %s (%s) : skipping layer\n", layer.Name, rtyq.ErrDatabaseFileDoesNotExist.Error(), fn)
			continue
		}

		fmt.Printf("initializing database\n")

		db, err := rtyq.InitDB(layer.Database.Path)
		if err != nil {
			fmt.Printf("warning : layer (%s) : %s (%s) : skipping layer\n", layer.Name, err.Error(), fn)
			continue
		}

		fmt.Printf("running spatial index:%s (%s)...\n", db.Index, db.FileName)
		start := time.Now()

		err = db.CreateSpatialIndex(layer.Database.Index)
		if err != nil {
			fmt.Printf("warning : layer (%s) : %s (%s) : skipping layer\n", layer.Name, err.Error(), fn)
			continue
		}

		dur := time.Since(start)
		fmt.Printf("time to index: %s\n", dur)

		handler := Handler{
			Data:      data,
			Database:  db,
			ZoomLimit: layer.Service.ZoomLimit,
		}

		layerEndpoint := fmt.Sprintf("/%s", layer.Service.Endpoint)

		layerEndpoints = append(layerEndpoints,
			fmt.Sprintf("%s/%s/%s", layerEndpoint, "id", "{id}"),
			fmt.Sprintf("%s/%s/%s", layerEndpoint, "point", "{lon,lat}"),
			fmt.Sprintf("%s/%s/%s", layerEndpoint, "tile", "{z/x/y}"),
		)

		router.Route(layerEndpoint, func(r chi.Router) {
			r.Get("/point/{point}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleLayer(w, r, "point", cfg.EnableLogs)
			})
			r.Get("/tile/{z}/{x}/{y}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleLayer(w, r, "tile", cfg.EnableLogs)
			})
			r.Get("/id/{id}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleLayer(w, r, "id", cfg.EnableLogs)
			})
		})

		numRoutes++
	}

	// write message to not found /

	notfound := Message{
		Message:   "invalid endpoint",
		Endpoints: layerEndpoints,
	}

	b, err := json.Marshal(notfound)
	if err != nil {
		return numRoutes, ErrUnableToWriteMessage
	}

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	return numRoutes, nil
}
