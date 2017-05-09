package kismet

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const websocketGUID string = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

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

// generateAcceptHeader returns the Base64-encoded SHA1 hash of the
// Sec-Websocket-Key header and the WebSocket GUID.
func (s *WebSocketServer) generateAcceptHeader(key string) string {
	hash := sha1.New()
	io.WriteString(hash, key)
	io.WriteString(hash, websocketGUID)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// upgradeConnection attempts to upgrade an incoming HTTP request to the
// WebSocket protocol.
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

	requestURL, _ := url.Parse(origin)

	// Requests for a WebSocket protocol upgrade should probably be disregarded
	// if they come from a different host.
	if !(s.config.AllowCrossOrigin == true || requestURL.Host == r.Host) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	connectionHeaders := strings.Split(r.Header.Get("Connection"), ", ")

	// Firefox sends keep-alive with the Connection header, so it's necessary to
	// check everything in the Connection header to make sure Upgrade is included.
	if !sliceContainsString(connectionHeaders, "Upgrade") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	conn, buffer, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	socket := NewWebSocket(conn, buffer)
	accept := s.generateAcceptHeader(r.Header.Get("Sec-Websocket-Key"))

	socket.buffer.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	socket.buffer.WriteString("Upgrade: websocket\r\n")
	socket.buffer.WriteString("Connection: Upgrade\r\n")
	socket.buffer.WriteString(fmt.Sprintf("Sec-WebSocket-Accept: %s\r\n", accept))
	socket.buffer.WriteString("\r\n\r\n")
	socket.buffer.Flush()

	defer socket.Close()
}
