package proxy

import (
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/monishth/dumb-prox/internal/middleware"
	"github.com/monishth/dumb-prox/internal/utils"
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

type ForwardProxy struct {
	client *http.Client
	once   sync.Once
}

func (p *ForwardProxy) gethttpClient() *http.Client {
	p.once.Do(func() {
		p.client = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
			// Probably don't want to chase redirects unless the client does
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: 10 * time.Second,
		}
	})
	return p.client
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
func (p *ForwardProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestID := middleware.GetRequestID(req.Context())
	log.Info("Request:", "ID", requestID, "Address", req.RemoteAddr, "Method", req.Method, "URL", req.URL, "Host", req.Host)
	if req.Method == http.MethodConnect {

	} else {
		p.handleHTTPRequest(w, req)
	}
}

func (p *ForwardProxy) handleHTTPRequest(w http.ResponseWriter, req *http.Request) {
	requestID := middleware.GetRequestID(req.Context())
	// Ignore none http/https
	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		msg := "unsupported protocol scheme " + req.URL.Scheme
		http.Error(w, msg, http.StatusBadRequest)
		log.Error(msg)
		return
	}

	req.RequestURI = ""

	// Append X-Forwarded-For
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		req.Header.Set("X-Forwarded-For", clientIP)
	}

	// remove hop by hop headers
	removeHopByHopHeaders(req.Header)

	resp, err := p.gethttpClient().Do(req)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()
	log.Info("Response:", "ID", requestID, "Address", req.RemoteAddr, "Status", resp.Status, "header", resp.Header)

	removeHopByHopHeaders(resp.Header)
	utils.CopyMapArray(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
