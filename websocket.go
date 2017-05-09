package kismet

import (
	"bufio"
	"net"
)

type WebSocket struct {
	conn   *net.Conn
	buffer *bufio.ReadWriter
}

// NewWebSocket returns a new WebSocket.
func NewWebSocket(conn *net.Conn, buffer *bufio.ReadWriter) *WebSocket {
	return &WebSocket{conn, buffer}
}
