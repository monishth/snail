package proxy

import (
	"bytes"
	"errors"
	"github.com/charmbracelet/log"

	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type netConnMock struct {
	readFunc  func(p []byte) (n int, err error)
	writeFunc func(b []byte) (n int, err error)
	closeFunc func() error
}

func (n *netConnMock) Read(p []byte) (int, error)         { return n.readFunc(p) }
func (n *netConnMock) Write(b []byte) (int, error)        { return n.writeFunc(b) }
func (n *netConnMock) Close() error                       { return n.closeFunc() }
func (n *netConnMock) LocalAddr() net.Addr                { return nil }
func (n *netConnMock) RemoteAddr() net.Addr               { return nil }
func (n *netConnMock) SetDeadline(t time.Time) error      { return nil }
func (n *netConnMock) SetReadDeadline(t time.Time) error  { return nil }
func (n *netConnMock) SetWriteDeadline(t time.Time) error { return nil }

func TestForwardProxy_handleHTTPSRequest(t *testing.T) {
	tests := []struct {
		name            string
		setupMockDialer func() func(network, address string) (net.Conn, error)
	}{
		{
			name: "successful dial",
			setupMockDialer: func() func(network, address string) (net.Conn, error) {
				return func(network, address string) (net.Conn, error) {
					return &netConnMock{
						readFunc: func(p []byte) (n int, err error) {
							return 0, nil
						},
						writeFunc: func(b []byte) (n int, err error) {
							return len(b), nil
						},
						closeFunc: func() error {
							return nil
						},
					}, nil
				}
			},
		},
		{
			name: "unsuccessful dial",
			setupMockDialer: func() func(network, address string) (net.Conn, error) {
				return func(network, address string) (net.Conn, error) {
					return nil, errors.New("mock dial error")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.New(&buf)
			req := httptest.NewRequest(http.MethodConnect, "https://example.com", nil)
			w := httptest.NewRecorder()
			p := &ForwardProxy{}
			http.DefaultTransport.(*http.Transport).Dial = tt.setupMockDialer()
			p.handleHTTPSRequest(w, req, logger)
		})
	}
}
