package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aliostad/TraceView/tracing"
)

type TraceApi struct {
	config *tracing.Config
	store  tracing.TraceStore
	server *http.Server
}

func NewTraceApi(port int,
	address string,
	config *tracing.Config,
	store tracing.TraceStore) *TraceApi {

	return &TraceApi{
		config: config,
		store:  store,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", address, port),
			Handler: nil, // to use default handler
		},
	}
}

func (api *TraceApi) Start() error {
	fs := http.FileServer(http.Dir("./content"))

	http.HandleFunc("/$", homePage)
	http.Handle("/", fs)

	go func() {
		if err := api.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func (api *TraceApi) Stop(ctx context.Context) error {
	return api.server.Shutdown(ctx)
}
