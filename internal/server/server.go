package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/monishth/snail/internal/auth"
	"github.com/monishth/snail/internal/middleware"
	"github.com/monishth/snail/internal/proxy"
)

func RunServer(options ServerOptions) {
	var addr = "127.0.0.1:" + strconv.Itoa(options.Port)
	var authProvider auth.AuthProvider
	switch options.AuthProvider {
	case "htpasswd":
		authProvider = auth.NewHtpasswdProvider(options.HttpasswdFilename)
	case "simple":
		creds := strings.SplitN(options.SimpleCredentials, ":", 2)
		if len(creds) != 2 {
			log.Fatal("invalid user/pass format. Use user:pass")
		}
		log.Info("Setting up simple auth")
		authProvider = &auth.SimpleAuthProvider{User: creds[0], Pass: creds[1]}
	}

	ctx, cancel := context.WithCancel(context.Background())
	proxy := &proxy.ForwardProxy{}

	go func() {
		log.Info("Starting proxy server on", "address", addr)
		var srv *http.Server

		if authProvider == nil {
			srv = &http.Server{Addr: addr, Handler: middleware.RequestLoggerMiddleware(proxy)}
		} else {
			srv = &http.Server{Addr: addr, Handler: middleware.RequestLoggerMiddleware(
				middleware.BasicAuthMiddleware(
					proxy,
					authProvider,
				))}
		}
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
