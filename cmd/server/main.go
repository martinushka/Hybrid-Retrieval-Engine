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
	"github.com/martinushka/ios-rag/internal/embedding"
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

	ctx := context.Background()

	conn, err := product.ConnectPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres connection error: %v", err)
	}
	defer conn.Close(ctx)

	repository := product.NewPostgresRepository(conn)
	embedder := embedding.NewHTTPProvider("http://127.0.0.1:8000")
	searchService := search.NewHybridService(repository, embedder)
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
		log.Printf("server started on :%s", cfg.Port)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop

	log.Println("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
