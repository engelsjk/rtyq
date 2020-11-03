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
func Start(port int, bdb *buntdb.DB, dirData, index string) error {

	router := chi.NewRouter()

	routes(router, bdb, dirData, index)

	fmt.Printf("running locally on :%d\n", port)

	return http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(port)), router)
}
