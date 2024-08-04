package proxy

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/monishth/snail/internal/utils"
)

func removeHopByHopHeaders(h http.Header) {
	for _, header := range hopHeaders {
		h.Del(header)
	}

	for _, header := range h["Connection"] {
		for _, option := range strings.Split(header, ",") {
			if option = strings.TrimSpace(option); option != "" {
				h.Del(option)
			}
		}
	}
}

func (p *ForwardProxy) handleHTTPRequest(w http.ResponseWriter, req *http.Request, logger *log.Logger) {
	removeHopByHopHeaders(req.Header)
	if req.URL.Scheme != "http" {
		msg := "unsupported protocol scheme " + req.URL.Scheme
		http.Error(w, msg, http.StatusBadRequest)
		logger.Error(msg)
		return
	}

	host := req.URL.Host
	if !strings.Contains(host, ":") {
		if req.URL.Scheme == "http" {
			host = host + ":80"
		}
	}

	serverConn, dialErr := net.Dial("tcp", host)
	if dialErr != nil {
		msg := "Failed to connect to upstream server"
		logger.Print("Error", "message", msg, "error", dialErr)
		http.Error(w, msg, http.StatusBadGateway)
		return
	}

	defer serverConn.Close()
	// Append X-Forwarded-For
	if clientIP, _, splitErr := net.SplitHostPort(req.RemoteAddr); splitErr == nil {
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		req.Header.Set("X-Forwarded-For", clientIP)
	}

	// Adjust the request to only include path + query
	req.RequestURI = ""
	req.URL.Scheme = ""
	req.URL.Host = ""

	if writeErr := req.Write(serverConn); writeErr != nil {
		msg := "Failed to write to upstream server"
		log.Error("Error", "Message", msg, "Host", host)
		http.Error(w, msg, http.StatusBadGateway)
	}

	resp, dialErr := http.ReadResponse(bufio.NewReader(serverConn), req)
	if dialErr != nil {
		msg := "Failed to read response from upstream server"
		logger.Error("Error", "Message", msg)
		log.Error("Error", "Message", msg, "Host", host)
		http.Error(w, msg, http.StatusBadGateway)
	}
	defer resp.Body.Close()

	// remove hop by hop headers
	removeHopByHopHeaders(resp.Header)

	logger.Info("Response:", "Address", req.RemoteAddr, "Status", resp.StatusCode)

	utils.CopyMapArray(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
