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

// Message will be served as the home endpoint to the service
type Message struct {
	Status    string   `json:"status"`
	Endpoints []string `json:"endpoints"`
}

// SetRoutes ...
func SetRoutes(router *chi.Mux, cfg *rtyq.Config) error {

	layerEndpoints := []string{} // initialize endpoint list

	// iterate over data layers (data/database/service), initialize db and set endpointss
	for _, layer := range cfg.Layers {

		fmt.Println("%************%")

		data, err := rtyq.InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
		if err != nil {
			return err
		}

		db, err := rtyq.LoadDB(layer.Database.Path)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), filepath.Base(layer.Database.Path))
		}

		err = db.CreateSpatialIndex(layer.Database.Index)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), filepath.Base(layer.Database.Path))
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
