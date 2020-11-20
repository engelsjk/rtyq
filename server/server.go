package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/engelsjk/rtyq/conf"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type serverError struct {
	Error   error
	Message string
	Code    int
}

type serverHandler func(http.ResponseWriter, *http.Request) *serverError

func Create() *http.Server {

	confServer := conf.Configuration.Server

	bindAddr := fmt.Sprintf("%v:%v", confServer.Host, confServer.Port)
	// log host:port and cors origin

	router := initRouter()

	// timeouts

	timeoutSecRequest := conf.Configuration.Server.WriteTimeoutSec
	timeoutSecWrite := timeoutSecRequest + 1

	// middleware

	corsOpt := cors.Options{
		AllowedOrigins:   []string{conf.Configuration.Server.CORSOrigin},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}

	router.Use(
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Compress(5, "gzip"),
		cors.Handler(corsOpt),
	)

	if conf.Configuration.Server.Logs {
		router.Use(middleware.Logger)
	}

	server := &http.Server{
		ReadTimeout:  time.Duration(conf.Configuration.Server.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(timeoutSecWrite) * time.Second,
		Addr:         bindAddr,
		Handler:      router,
	}

	return server
}

func (sh serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	handlerDone := make(chan struct{})
	start := time.Now()
	_ = start

	go func() {
		select {
		case <-handlerDone:
			// log time to complete request
		case <-r.Context().Done():
			// log canceled requests
			switch r.Context().Err() {
			case context.DeadlineExceeded:
				// log request terminated by write timeout
			case context.Canceled:
				// log request canceled by client
			}
		}
	}()

	e := sh(w, r)

	if e != nil {
		// log request processing error
		http.Error(w, e.Message, e.Code)
	}
	close(handlerDone)
}

func serverErrorInternal(err error, msg string) *serverError {
	return &serverError{err, msg, http.StatusInternalServerError}
}

func FatalAfter(delaySec int, msg string) chan struct{} {
	chanCancel := make(chan struct{})
	go func() {
		select {
		case <-chanCancel:
			return
		case <-time.After(time.Duration(delaySec) * time.Second):
			// log msg
		}
	}()
	return chanCancel
}
