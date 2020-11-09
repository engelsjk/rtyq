package service

import (
	"fmt"
	"net/http"

	"github.com/engelsjk/rtyq/pkg/config"
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/go-chi/chi"
	"github.com/tidwall/buntdb"
)

// routes ...
func routes(router *chi.Mux, cfg *config.Config) error {

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("query home"))
	})

	bdbs := []*buntdb.DB{}

	for i, set := range cfg.Sets {

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

		router.Get(endpoint, func(w http.ResponseWriter, r *http.Request) {
			handler.HandleData(w, r)
		})

		fmt.Printf("endpoint %s set\n", endpoint)
	}

	return nil
}
