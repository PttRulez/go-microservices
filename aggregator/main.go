package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pttrulez/go-microservices/aggregator/client"
	"github.com/pttrulez/go-microservices/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpaddr", ":3000", "The address to listen on for HTTP requests.")
	grpcListenAddr := flag.String("grpcaddr", ":3001", "The address to listen on for GRPC requests.")
	flag.Parse()

	var (
		store = NewMemortyStore()
		svc   = NewLogMiddleware(NewInvoiceAggregator(store))
	)
	go makeGRPCTransport(*grpcListenAddr, svc)
	// go makeHTTPTransport(*httpListenAddr, svc)
	time.Sleep(time.Second * 5)
	c, err := client.NewGRPCClient(*grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GRPC Client created")
	if _, err = c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 58.55,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
	log.Fatal(makeHTTPTransport(*httpListenAddr, svc))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	// Make a TCP Listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	// Make a new GRPC native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)

	// Register (OUR) GRPC server implementation to the GRPC package.
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	fmt.Println("GRPC Transport is running on port", listenAddr)
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP Transport is running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
	return http.ListenAndServe(listenAddr, nil)
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
