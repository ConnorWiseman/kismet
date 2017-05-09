package kismet

import (
	"net/http"
)

type WebSocketServer struct {
}

// NewWebSocketServer returns a new WebSocketServer.
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{}
}

// ServeHTTP ensures that WebSocketServer fulfills the http.Handler interface.
func (s *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
