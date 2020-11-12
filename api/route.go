package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/engelsjk/rtyq"
	"github.com/go-chi/chi"
)

var (
	ErrUnableToWriteHomeMessage error = fmt.Errorf("unable to write home message")
)

// Message acts as the home endpoint output for the api service
type Message struct {
	Status    string   `json:"status"`
	Endpoints []string `json:"endpoints"`
}

// SetRoutes initializes all of the API service endpoints. 	/
// It iterates over each layer to initialize each db with a spatial index
// and links it to a separate layer endpoint.
func SetRoutes(router *chi.Mux, cfg *rtyq.Config) error {

	layerEndpoints := []string{}

	for _, layer := range cfg.Layers {

		fn := filepath.Base(layer.Database.Path)

		fmt.Println("%************%")

		data, err := rtyq.InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), fn)
		}

		db, err := rtyq.InitDB(layer.Database.Path)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), fn)
		}

		err = db.CreateSpatialIndex(layer.Database.Index)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), fn)
		}

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
		Status:    "running",
		Endpoints: layerEndpoints,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return ErrUnableToWriteHomeMessage
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})

	return nil
}
