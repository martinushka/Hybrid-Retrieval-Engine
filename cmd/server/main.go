package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/martinushka/ios-rag/internal/api"
	"github.com/martinushka/ios-rag/internal/config"
	"github.com/martinushka/ios-rag/internal/httpserver"
	"github.com/martinushka/ios-rag/internal/product"
	"github.com/martinushka/ios-rag/internal/search"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	cfg := config.Load()
	products := []product.Product{
		{
			ID:          "1",
			Title:       "Lenovo ThinkPad X1",
			Description: "Lightweight business laptop for professionals",
			Category:    "laptops",
			Price:       1299.99,
		},
		{
			ID:          "2",
			Title:       "MacBook Air M2",
			Description: "Thin and lightweight laptop for everyday work",
			Category:    "laptops",
			Price:       999.99,
		},
		{
			ID:          "3",
			Title:       "ASUS Zenbook 14",
			Description: "Compact laptop with OLED display",
			Category:    "laptops",
			Price:       1099.99,
		},
	}

	searchService := search.NewInMemoryService(products)
	apiHandler := api.NewHandler(searchService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/search", apiHandler.Search)

	handler := httpserver.Logging(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("ios-rag server started on :%s", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop

	log.Println("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
