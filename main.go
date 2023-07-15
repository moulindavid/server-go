package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	cfg := &apiConfig{fileserverHits: 0}
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", serveHealth)
	mux.HandleFunc("/metrics", cfg.serveMetrics)
	corsMux := middlewareCors(mux)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMux))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func serveHealth(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) serveMetrics(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
