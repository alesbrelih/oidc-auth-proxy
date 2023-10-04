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

	"github.com/alesbrelih/go-oidc-auth-proxy/cmd/go-oidc-auth-proxy/handler"
	"github.com/alesbrelih/go-oidc-auth-proxy/internal/config"
	"github.com/alesbrelih/go-oidc-auth-proxy/internal/generated/oidc/api"
	oidcPkg "github.com/alesbrelih/go-oidc-auth-proxy/internal/oidc"
	"github.com/alesbrelih/go-oidc-auth-proxy/internal/transform"
	_ "github.com/hashicorp/cap"
)

func main() {
	baseCtx := context.TODO()

	cfg := config.Config{}
	err := config.Get(&cfg)
	if err != nil {
		log.Fatalf("could not parse config: %s", err)
	}

	oidcSvc, err := oidcPkg.New(baseCtx, cfg)
	if err != nil {
		log.Fatalf("cloud not initialize OIDC provider: %s", err)
	}

	t, err := transform.New(transform.DefaultTemplate)
	if err != nil {
		log.Fatalf("could not create HTTP header value transformer: %s", err)
	}

	log.Printf("cfg: %+v", cfg)

	handler, err := api.NewServer(handler.New(oidcSvc, t))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: LoggingMiddleware(handler),
	}

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Server started @%d", cfg.Port)

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-osSignal

	contextTimeout := 5 * time.Second
	ctxShutdown, cancelCtxShutdown := context.WithTimeout(context.Background(), contextTimeout)
	defer cancelCtxShutdown()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("error gracefully shutting down server: %s", err)
	}

}

// LoggingMiddleware logs some basic info about each HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the received request
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Record the start time
		startTime := time.Now()

		// Call the next middleware/handler in the chain
		next.ServeHTTP(w, r)

		// Log the completed request
		log.Printf(
			"Completed %s %s in %v",
			r.Method,
			r.URL.Path,
			time.Since(startTime),
		)
	})
}
