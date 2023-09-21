package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func handleSingle(w http.ResponseWriter, r *http.Request) {
	sleepTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepTime)
	fmt.Fprintf(w, "slept %v\n", sleepTime)
}

func handleMulti(w http.ResponseWriter, r *http.Request) {
	subRequests := 3 + rand.Intn(4)
	slog.InfoContext(r.Context(), "in /long request handler", slog.Int("subRequests", subRequests))

	for i := 0; i < subRequests; i++ {
		if err := callSingle(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}

	fmt.Fprintln(w, "ok")
}

func runServer() error {
	http.Handle("/single", otelhttp.NewHandler(http.HandlerFunc(handleSingle), "/single"))
	http.Handle("/multi", otelhttp.NewHandler(http.HandlerFunc(handleMulti), "/multi"))
	return http.ListenAndServe(":8080", nil)
}
