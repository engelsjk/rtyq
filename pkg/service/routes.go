package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/engelsjk/rtyq/pkg/config"
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/go-chi/chi"
	"github.com/tidwall/buntdb"
)

// Message will be served as the home endpoint to the service
type Message struct {
	Status    string   `json:"status"`
	Endpoints []string `json:"endpoints"`
}

// routes ...
func setRoutes(router *chi.Mux, cfg *config.Config) error {

	var (
		ErrUnableToWriteHomeMessage error = fmt.Errorf("unable to write home message")
	)

	endpoints := []string{} // initialize endpoint list
	bdbs := []*buntdb.DB{}  // initialize slice of db pointers

	// iterate over data sets (data/database/service), initialize db and set endpoint
	for i, set := range cfg.Sets {

		fmt.Println("%************%")

		bdbi, err := db.Initialize(set.Database.Path, set.Database.Index, false)
		if err != nil {
			return err
		}

		bdbs = append(bdbs, bdbi)

		handler := Handler{
			DirData:   set.Data.Path,
			Extension: set.Data.Extension,
			Database:  bdbs[i],
			Index:     set.Database.Index,
			ZoomLimit: set.Service.ZoomLimit,
		}

		endpoint := fmt.Sprintf("/%s", set.Service.Path)

		endpoints = append(endpoints, endpoint)

		router.Get(endpoint, func(w http.ResponseWriter, r *http.Request) {
			handler.HandleData(w, r)
		})
	}

	// write message to home /

	message := Message{
		Status:    "running",
		Endpoints: endpoints,
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
