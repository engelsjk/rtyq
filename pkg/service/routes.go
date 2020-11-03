package service

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tidwall/buntdb"
)

func routes(router *chi.Mux, dirData, ext string, bdb *buntdb.DB, index string) {

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("query home"))
	})

	endpoint := fmt.Sprintf("/%s", index)

	router.Get(endpoint, func(w http.ResponseWriter, r *http.Request) {
		HandleData(w, r, dirData, ext, bdb, index)
	})
}
