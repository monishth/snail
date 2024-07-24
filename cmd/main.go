package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/monishth/dumb-prox/internal/auth"
	"github.com/monishth/dumb-prox/internal/middleware"
	"github.com/monishth/dumb-prox/internal/proxy"
)

func main() {
	var addr = flag.String("addr", "127.0.0.1:9999", "proxy address")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	proxy := &proxy.ForwardProxy{}

	go func() {
		log.Info("Starting proxy server on", "address", *addr)
		htpasswd := auth.NewHtpasswdProvider("passwd_file")

		srv := &http.Server{Addr: *addr, Handler: middleware.RequestLoggerMiddleware(
			middleware.BasicAuthMiddleware(
				proxy,
				// &auth.SimpleAuthProvider{User: "test", Pass: "test"},
				&htpasswd,
			))}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("ListenAndServe:", err)
			}
		}()

		<-ctx.Done()
		ctxShutDown, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	log.Info("Shutting down gracefully")
}
