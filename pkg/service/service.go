package service

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/tidwall/buntdb"
)

// Start ...
func Start(port int, dirData, ext string, bdb *buntdb.DB, index string) error {

	router := chi.NewRouter()

	routes(router, dirData, ext, bdb, index)

	fmt.Printf("running locally on :%d\n", port)

	return http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(port)), router)
}
