package app

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func handlePing() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Pong")
		},
	)
}

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/ping", handlePing())
}

func NewServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
	)
	var handler http.Handler = mux
	return handler
}

type Config struct {
	Host string
	Port string
}

func NewConfig(host string, port string) *Config {
	cfg := Config{
		Host: host,
		Port: port,
	}
	return &cfg
}

func Run(ctx context.Context, cfg *Config, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	srv := NewServer()
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: srv,
	}
	go func() {
		fmt.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
