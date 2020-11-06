package service

import (
	"fmt"
	"net/http"

	"github.com/engelsjk/rtyq/pkg/config"
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/go-chi/chi"
	"github.com/tidwall/buntdb"
)

func routes(router *chi.Mux, cfg *config.Config) error {

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("query home"))
	})

	bdbs := []*buntdb.DB{}

	for i, svc := range cfg.Services {

		bdbi, err := db.Initialize(svc.Database.Path, svc.Database.Index, false)
		if err != nil {
			return err
		}

		bdbs = append(bdbs, bdbi)

		endpoint := fmt.Sprintf("/%s", svc.Database.Index)

		router.Get(endpoint, func(w http.ResponseWriter, r *http.Request) {
			HandleData(w, r, svc.Data.Path, svc.Data.Extension, bdbs[i], svc.Database.Index)
		})

		fmt.Printf("endpoint %s set\n", endpoint)
	}

	return nil
}
