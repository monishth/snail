package proxy

import (
	"github.com/monishth/dumb-prox/internal/middleware"
	"net/http"
)

var hopHeaders = []string{
	"Connection",
	"Proxy-Authorization",
	"Proxy-Connection",
	"Proxy-Authenticate",
	"Keep-Alive",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

type ForwardProxy struct {
}

func (p *ForwardProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := middleware.GetLogger(req.Context())
	logger.Info("Request:", "Address", req.RemoteAddr, "Method", req.Method, "URL", req.URL, "Host", req.Host)

	if req.Method == http.MethodConnect {
		p.handleHTTPSRequest(w, req, logger)
	} else {
		p.handleHTTPRequest(w, req, logger)
	}
}
