package service

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/engelsjk/rtyq/pkg/config"
	"github.com/go-chi/chi"
)

// Start ...
func Start(cfg *config.Config) error {

	router := chi.NewRouter()

	err := routes(router, cfg)
	if err != nil {
		return err
	}

	fmt.Println("%************%")

	fmt.Printf("running %d services locally on localhost:%d\n", len(cfg.Sets), cfg.Port)

	return http.ListenAndServe(net.JoinHostPort("", strconv.Itoa(cfg.Port)), router)
}
