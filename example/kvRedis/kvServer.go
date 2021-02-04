package kvserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"

	quic "github.com/lucas-clemente/quic-go"
)

const h09alpn = "quic-echo-example"

// Server is a custom server listening for QUIC connections.
type Server struct {
	QuicConfig *quic.Config
	Addr       string
	TLSConfig  *tls.Config
	mutex      sync.Mutex
	listener   quic.EarlyListener
}

// Close closes the server.
func (s *Server) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.listener.Close()
}

func (s *Server) ListenAndServe() error {

	udpAddr, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	tlsConf := s.TLSConfig.Clone()
	tlsConf.NextProtos = []string{h09alpn}
	ln, err := quic.ListenEarly(conn, tlsConf, s.QuicConfig)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	s.listener = ln
	s.mutex.Unlock()

	for {
		sess, err := ln.Accept(context.Background())
		if err != nil {
			return err
		}
		go s.handleConn(sess)
	}
}

func (s *Server) handleConn(sess quic.Session) {
	for {
		str, err := sess.AcceptStream(context.Background())
		if err != nil {
			log.Printf("Error accepting stream: %s\n", err.Error())
			return
		}
		go func() {
			if err := s.handleStream(str); err != nil {
				log.Printf("Handling stream failed: %s\n", err.Error())
			}
		}()
	}
}

func (s *Server) handleStream(stream quic.Stream) error {
	fmt.Println("to handleStream")
	// _, err = io.Copy(loggingWriter{stream}, stream)

	reqBytes, err := ioutil.ReadAll(stream)
	if err != nil {
		return err
	}
	request := string(reqBytes)
	request = strings.TrimRight(request, "\r\n")
	// request = strings.TrimRight(request, " ")

	log.Printf("Received request: %s\n", request)
	return nil
}
