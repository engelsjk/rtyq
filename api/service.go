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

	err = SetRoutes(router, cfg)
	if err != nil {
		return err
	}

	plural := ""
	if len(cfg.Layers) > 1 {
		plural = "s"
	}

	fmt.Println("%************%")
	fmt.Printf("starting %d layer service%s on localhost:%d\n", len(cfg.Layers), plural, cfg.Port)

	return http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(cfg.Port)), router)
}

// SetRoutes initializes all of the API service endpoints. 	/
// It iterates over each layer to initialize each db with a spatial index
// and links it to a separate layer endpoint.
func SetRoutes(router *chi.Mux, cfg *rtyq.Config) error {

	layerEndpoints := []string{}

	for _, layer := range cfg.Layers {

		fn := filepath.Base(layer.Database.Path)

		fmt.Println("%************%")
		fmt.Printf("setting route for layer: %s\n", layer.Name)

		data, err := rtyq.InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), fn)
		}

		fmt.Printf("initializing database\n")

		db, err := rtyq.InitDB(layer.Database.Path)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), fn)
		}

		fmt.Printf("running spatial index:%s (%s)...\n", db.Index, db.FileName)
		start := time.Now()

		err = db.CreateSpatialIndex(layer.Database.Index)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), fn)
		}

		dur := time.Since(start)
		fmt.Printf("time to index: %s\n", dur)

		handler := Handler{
			Data:      data,
			Database:  db,
			ZoomLimit: layer.Service.ZoomLimit,
		}

		layerEndpoint := fmt.Sprintf("/%s", layer.Service.Endpoint)

		layerEndpoints = append(layerEndpoints, layerEndpoint)

		router.Route(layerEndpoint, func(r chi.Router) {
			r.Get("/point/{point}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleLayer(w, r, "point")
			})
			r.Get("/tile/{z}/{x}/{y}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleLayer(w, r, "tile")
			})
			r.Get("/id/{id}", func(w http.ResponseWriter, r *http.Request) {
				handler.HandleLayer(w, r, "id")
			})
		})
	}

	// write message to home /

	message := Message{
		Message:   "running",
		Endpoints: layerEndpoints,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return ErrUnableToWriteMessage
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	// write message to home /

	message = Message{
		Message:   "endpoint not found",
		Endpoints: layerEndpoints,
	}

	b, err = json.Marshal(message)
	if err != nil {
		return ErrUnableToWriteMessage
	}

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	return nil
}
