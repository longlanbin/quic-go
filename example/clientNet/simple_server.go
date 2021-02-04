package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"

	quic "github.com/lucas-clemente/quic-go"
)

const host = "172.20.114.49"
const port = 9981

func handleSession(sess quic.Session) {
	for {
		fmt.Println("waiting for coming data")
		stream, err := sess.AcceptStream(context.Background())
		fmt.Println("=============================1")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		// Echo through the loggingWriter
		fmt.Println("to copy coming data")
		_, err = io.Copy(loggingWriter{stream}, stream)
	}
}

func main() {
	addr := host + ":" + strconv.Itoa(port)
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("waiting for coming session")
		sess, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleSession(sess)
	}
	// go func() { echoServer() }()

	// listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	// data := make([]byte, 1024)
	// for {
	// 	n, remoteAddr, err := listener.ReadFromUDP(data)
	// 	if err != nil {
	// 		fmt.Printf("error during read: %s", err)
	// 	}
	// 	fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
	// 	_, err = listener.WriteToUDP([]byte("world"), remoteAddr)
	// 	if err != nil {
	// 		fmt.Printf(err.Error())
	// 	}
	// }

}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {
	addr := host + ":" + strconv.Itoa(port)
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	sess, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}
	// Echo through the loggingWriter
	_, err = io.Copy(loggingWriter{stream}, stream)
	return err
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	keyLog, err := os.Create("serverkey.log")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("create key log: %s\n", keyLog)
	if err != nil {
		fmt.Printf("Could not create key log: %s\n", err.Error())
		os.Exit(1)
	}
	if keyLog != nil {
		defer keyLog.Close()
	}
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
		// KeyLogWriter: keyLog,
	}
}
