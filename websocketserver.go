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
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	origin := r.Header.Get("Origin")

	// Requests without an Origin header should be denied.
	if origin == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	url, _ := url.Parse(origin)

	// Requests for a WebSocket protocol upgrade should probably be disregarded
	// if they come from a different host.
	if !(s.config.AllowCrossOrigin == true || r.Host == r.Host) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	connectionHeaders := strings.Split(r.Header.Get("Connection"), ", ")

	// Firefox sends keep-alive with the Connection header, so it's necessary to
	// check everything in the Connection headers to make sure Upgrade is included.
	if !sliceContainsString(connectionHeaders, "Upgrade") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

  conn, buffer, err := w.(http.Hijacker).Hijack()
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  defer conn.Close()
}
