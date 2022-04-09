package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aliostad/TraceView/tracing"
)

var (
	singletonApi *TraceApi
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

	if singletonApi != nil {
		panic("Another TraceAPI already exists.")
	}

	singletonApi = &TraceApi{
		config: config,
		store:  store,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", address, port),
			Handler: nil, // to use default handler
		},
	}

	return singletonApi
}

func (api *TraceApi) Start() error {
	fs := http.FileServer(http.Dir("./content"))

	http.HandleFunc("/api/traces", traces)
	http.HandleFunc("/$", homePage)
	http.Handle("/", fs)

	go func() {
		if err := api.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func traces(w http.ResponseWriter, r *http.Request) {
	var to, from *time.Time
	var n int
	froms := r.URL.Query().Get("from")
	tos := r.URL.Query().Get("to")
	counts := r.URL.Query().Get("count")

	if froms == "" {
		from = nil
	} else {
		fromX, err := time.Parse(time.RFC3339, froms)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		from = &fromX
	}

	if tos == "" {
		toX := time.Now().UTC()
		to = &toX
	} else {
		toX, err := time.Parse(time.RFC3339, froms)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		to = &toX
	}

	if counts == "" {
		n = 100
	} else {
		n, _ = strconv.Atoi(counts)
	}

	traces, err := singletonApi.store.ListByTimeRange(n, from, to)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(traces)
}

func (api *TraceApi) Stop(ctx context.Context) error {
	return api.server.Shutdown(ctx)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusFound)
}
