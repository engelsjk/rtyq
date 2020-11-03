package service

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tidwall/buntdb"
)

func routes(router *chi.Mux, bdb *buntdb.DB, dirData, index string) {

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("query home"))
	})

	router.Get("/block", func(w http.ResponseWriter, r *http.Request) {
		HandleData(w, r, bdb, dirData, index)
	})
}
