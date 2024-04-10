package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strconv"

	"github.com/pttrulez/toll-calc/types"
	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "The address to listen on for HTTP requests.")
	flag.Parse()

	var (
		store = NewMemortyStore()
		svc   = NewLogMiddleware(NewInvoiceAggregator(store))
	)
	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuID := r.URL.Query().Get("obuid")
		if obuID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "obuid is required"})
		} else {
			obuIDint, err := strconv.Atoi(obuID)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			}

			invoice, err := svc.CalcualateInvoice(obuIDint)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, invoice)
		}
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			logrus.Error("json decode error: ", err)
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			logrus.Error("aggregation error: ", err)
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
