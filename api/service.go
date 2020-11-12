package api

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/engelsjk/rtyq"
	"github.com/go-chi/chi"
)

// StartService ...
func StartService(cfg *rtyq.Config) error {

	if cfg == nil {
		return fmt.Errorf("no config provided")
	}

	router := chi.NewRouter()

	err := SetRoutes(router, cfg)
	if err != nil {
		return err
	}

	fmt.Println("%************%")

	fmt.Printf("starting %d layer services locally on localhost:%d\n", len(cfg.Layers), cfg.Port)

	return http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(cfg.Port)), router)
}
