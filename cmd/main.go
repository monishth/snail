package main

import (
	"flag"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/monishth/dumb-prox/internal/auth"
	"github.com/monishth/dumb-prox/internal/middleware"
	"github.com/monishth/dumb-prox/internal/proxy"
)

func main() {
	var addr = flag.String("addr", "127.0.0.1:9999", "proxy address")
	flag.Parse()

	proxy := &proxy.ForwardProxy{}

	log.Info("Starting proxy server on", "address", *addr)
	htpasswd := auth.NewHtpasswdProvider("passwd_file")
	if err := http.ListenAndServe(
		*addr,
		middleware.RequestLoggerMiddleware(
			middleware.BasicAuthMiddleware(
				proxy,
				// &auth.SimpleAuthProvider{User: "test", Pass: "test"},
				&htpasswd,
			)),
	); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
