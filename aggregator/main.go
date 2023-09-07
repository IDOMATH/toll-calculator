package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/idomath/toll-calculator/types"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strconv"
)

func main() {
	httpListenPort := flag.String("httpListenPort", ":3000", "listen port for the HTTP server")
	grpcListenPort := flag.String("grpcListenPort", ":3001", "listen port for the gRPC server")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGrpcTransport(*grpcListenPort, svc))
	}()
	log.Fatal(makeHttpTransport(*httpListenPort, svc))
}

func makeHttpTransport(listenPort string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port: ", listenPort)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice(svc))
	return http.ListenAndServe(listenPort, nil)
}

func makeGrpcTransport(listenPort string, svc Aggregator) error {
	fmt.Println("gRPC transport running on port: ", listenPort)
	// Make a TCP listener
	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		return err
	}
	defer listener.Close()

	// Make a new gRPC native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)

	// Register our gRPC server implementation with the gRPC package.
	types.RegisterAggregatorServer(server, NewGrpcAggregatorServer(svc))
	return server.Serve(listener)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obuId"]
		if !ok {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuId)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, http.StatusOK, invoice)
		return
	}
}

func writeJson(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
