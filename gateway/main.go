package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/idomath/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenPort := flag.String("listenPort", ":6000", "listen port of the HTTP server")
	aggregatorServiceAddress := flag.String("aggServiceAddr", "http://127.0.0.1:3000", "address of the aggregator service")

	client := client.NewHttpClient(*aggregatorServiceAddress) // Aggregator service endpoint
	invoiceHandler := NewInvoiceHandler(client)
	http.HandleFunc("/invoice", makeApiFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("gateway HTTP server running on port %s", *listenPort)
	log.Fatal(http.ListenAndServe(*listenPort, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := h.client.GetInvoice(context.Background(), 175732606)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, inv)
}

func writeJson(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info()
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
