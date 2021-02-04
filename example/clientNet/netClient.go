package main

import (
	// "bufio"
	// "bytes"
	// "crypto/tls"
	// "crypto/x509"
	// "flag"
	"fmt"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net"
	"os"
	"time"
	// "io"
	// "log"
	// "net/http"
	// "os"
	// "sync"
	// "github.com/lucas-clemente/quic-go"
	// "github.com/lucas-clemente/quic-go/http3"
	// "github.com/lucas-clemente/quic-go/internal/testdata"
	// "github.com/lucas-clemente/quic-go/internal/utils"
	// "github.com/lucas-clemente/quic-go/logging"
	// "github.com/lucas-clemente/quic-go/qlog"
)

func main() {
	fmt.Printf("hello longlan\n")

	var (
		cl              *client
		packetConn      *mockPacketConn
		addr            net.Addr
		connID          protocol.ConnectionID
		mockMultiplexer *MockMultiplexer
		origMultiplexer multiplexer
		tlsConf         *tls.Config
		tracer          *mocklogging.MockConnectionTracer
		config          *Config

		originalClientSessConstructor func(
			conn sendConn,
			runner sessionRunner,
			destConnID protocol.ConnectionID,
			srcConnID protocol.ConnectionID,
			conf *Config,
			tlsConf *tls.Config,
			initialPacketNumber protocol.PacketNumber,
			initialVersion protocol.VersionNumber,
			enable0RTT bool,
			hasNegotiatedVersion bool,
			tracer logging.ConnectionTracer,
			logger utils.Logger,
			v protocol.VersionNumber,
		) quicSession
	)

	// generate a packet sent by the server that accepts the QUIC version suggested by the client
	acceptClientVersionPacket := func(connID protocol.ConnectionID) []byte {
		b := &bytes.Buffer{}
		Expect((&wire.ExtendedHeader{
			Header:          wire.Header{DestConnectionID: connID},
			PacketNumber:    1,
			PacketNumberLen: 1,
		}).Write(b, protocol.VersionWhatever)).To(Succeed())
		return b.Bytes()
	}




	manager := NewMockPacketHandlerManager(mockCtrl)
	manager.EXPECT().Add(gomock.Any(), gomock.Any())
	manager.EXPECT().Destroy()
	mockMultiplexer.EXPECT().AddConn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(manager, nil)

	remoteAddrChan := make(chan string, 1)
	newClientSession = func(
		conn sendConn,
		_ sessionRunner,
		_ protocol.ConnectionID,
		_ protocol.ConnectionID,
		_ *Config,
		_ *tls.Config,
		_ protocol.PacketNumber,
		_ protocol.VersionNumber,
		_ bool,
		_ bool,
		_ logging.ConnectionTracer,
		_ utils.Logger,
		_ protocol.VersionNumber,
	) quicSession {
		remoteAddrChan <- conn.RemoteAddr().String()
		sess := NewMockQuicSession(mockCtrl)
		sess.EXPECT().run()
		sess.EXPECT().HandshakeComplete().Return(context.Background())
		return sess
	}
	_, err := DialAddr("localhost:17890", tlsConf, &Config{HandshakeTimeout: time.Millisecond})
}
