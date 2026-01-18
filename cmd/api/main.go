package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iramkrish/search-bff/internal/infra"

	internalhttp "github.com/iramkrish/search-bff/internal/http"
)

// main starts an HTTP server on :8080 that serves the /search endpoint and shuts down gracefully on SIGINT or SIGTERM.
// The server applies middleware for request IDs and logging, enforces read/write/idle timeouts, and allows up to 5 seconds for graceful shutdown.
func main() {
	logger := infra.NewLogger()

	mux := http.NewServeMux()
	handler := internalhttp.NewHandler(logger)

	mux.Handle("/search", handler)

	server := &http.Server{
		Addr: ":8080",
		Handler: internalhttp.Chain(
			mux,
			internalhttp.RequestID(),
			// internalhttp.Timeout(2*time.Second),
			internalhttp.Logging(logger),
		),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.Println("server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	logger.Println("shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("shutdown failed: %v", err)
	}

	logger.Println("server stopped")
}