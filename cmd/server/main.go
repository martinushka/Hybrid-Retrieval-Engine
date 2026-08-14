package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/martinushka/ios-rag/internal/config"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	cfg := config.Load()

	http.HandleFunc("/health", healthHandler)

	log.Printf("ios-rag server started on :%s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
