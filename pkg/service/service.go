package service

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/engelsjk/rtyq/pkg/config"
	"github.com/go-chi/chi"
)

// Service ...
var Service struct {
	Status  string
	Message map[string]interface{}
}

// Start ...
func Start(cfg *config.Config) error {

	router := chi.NewRouter()

	err := setRoutes(router, cfg)
	if err != nil {
		return err
	}

	fmt.Println("%************%")

	fmt.Printf("starting %d services locally on localhost:%d\n", len(cfg.Sets), cfg.Port)

	return http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(cfg.Port)), router)
}
