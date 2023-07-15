package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	cfg := apiConfig{fileserverHits: 0}

	r := chi.NewRouter()

	router_api := chi.NewRouter()
	handler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", handler)
	r.Handle("/app/*", handler)

	router_api.Get("/healthz", serveHealth)
	router_api.Get("/metrics", cfg.serveMetrics)

	r.Mount("/api", router_api)

	corsMux := middlewareCors(r)
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
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
