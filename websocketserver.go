package kismet

import (
	"net/http"
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
	return s.upgradeConnection(w, r)
}

// upgradeConnection attempts to upgrade an incoming HTTP request to the WebSocket protocol.
func (s *WebSocketServer) upgradeConnection(w http.ResponseWriter, r *http.Request) {

}
