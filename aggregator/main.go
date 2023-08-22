package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/idomath/toll-calculator/types"
	"net/http"
)

func main() {
	listenPort := flag.String("listenPort", ":3000", "listen port for the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	var (
		svc = NewInvoiceAggregator(store)
	)

	makeHttpTransport(*listenPort, svc)
}

func makeHttpTransport(listenPort string, svc Aggregator) {
	fmt.Println("HTTP transport running on port: ", listenPort)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenPort, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
