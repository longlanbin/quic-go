package kvserver

import (
	"crypto/tls"
	"sync"

	quic "github.com/lucas-clemente/quic-go"
)

// Server is a custom server listening for QUIC connections.
type Client struct {
	QuicConfig *quic.Config
	Addr       string
	TLSConfig  *tls.Config
	mutex      sync.Mutex
	listener   quic.EarlyListener
}
