package main

import (
	"ascii-art-web/handlers"
	"ascii-art-web/web/static"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	staticFS := http.FS(static.StaticFS)
	fs := http.FileServer(staticFS)
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.HandleHome)
	mux.HandleFunc("/ascii-art", handlers.HandleAsciiArt)

	fmt.Println("Starting server on :8080")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  responseTimeout,
		WriteTimeout: responseTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // block here until signal received

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}

var responseTimeout = 5 * time.Second
