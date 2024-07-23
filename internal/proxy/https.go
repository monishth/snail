package proxy

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
)

func (p *ForwardProxy) handleHTTPSRequest(w http.ResponseWriter, req *http.Request, logger *log.Logger) {
	host := req.URL.Host
	if !strings.Contains(host, ":") {
		if req.URL.Scheme == "https" {
			host = host + ":443"
		}
	}

	// Open TCP connection
	destConn, err := net.Dial("tcp", host)
	if err != nil {
		logger.Error("Error:", "Error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	// Turn response writer to hijacker, this lets us take over the raw connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		logger.Error("Error:", "message", "Hijacking not supported")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		logger.Error("Error", "message", "Error hijacking connection", "Error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
	logger.Info("Tunnel Success", "Address", req.RemoteAddr)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
