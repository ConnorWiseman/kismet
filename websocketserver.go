package kismet

import (
	"net/http"
	"net/url"
)

type WebSocketServer struct {
	config Config
}

// NewWebSocketServer returns a new WebSocketServer.
func NewWebSocketServer(config Config) *WebSocketServer {
	return &WebSocketServer{config}
}

// ServeHTTP ensures that WebSocketServer fulfills the http.Handler interface.
func (s *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.upgradeConnection(w, r)
	return
}

// upgradeConnection attempts to upgrade an incoming HTTP request to the WebSocket protocol.
func (s *WebSocketServer) upgradeConnection(w http.ResponseWriter, r *http.Request) {

	// The WebSocket handshake is only supported by HTTP/1.1 or greater.
	if !r.ProtoAtLeast(1, 1) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// The WebSocket handshake is only supported via GET.
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	origin := req.Header.Get("Origin")

	// Requests without an Origin header should be denied.
	if origin == "" {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
}
