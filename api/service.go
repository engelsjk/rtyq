package api

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/engelsjk/rtyq"
	"github.com/go-chi/chi"
)

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

	fmt.Println("%************%")

	fmt.Printf("starting %d layer services locally on localhost:%d\n", len(cfg.Layers), cfg.Port)

	return http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(cfg.Port)), router)
}
