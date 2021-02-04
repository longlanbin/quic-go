package main

import(
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/example/kvRedis"
	// "github.com/lucas-clemente/quic-go/interop/utils"
	// "github.com/lucas-clemente/quic-go/qlog"
)
const host = "0.0.0.0:9982"

func main(){
	logFile, err := os.Create("kvserverlog.txt")
	if err != nil {
		fmt.Printf("Could not create log file: %s\n", err.Error())
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// keyLog, err := utils.GetSSLKeyLog()
	// if err != nil {
	// 	fmt.Printf("Could not create key log: %s\n", err.Error())
	// 	os.Exit(1)
	// }
	// if keyLog != nil {
	// 	defer keyLog.Close()
	// }

	// getLogWriter, err := utils.GetQLOGWriter()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }
	// a quic.Config that doesn't do a Retry
	quicConf := &quic.Config{
		AcceptToken: func(_ net.Addr, _ *quic.Token) bool { return true },
		// Tracer:      qlog.NewTracer(getLogWriter),
	}
	cert, err := tls.LoadX509KeyPair("cert.pem", "priv.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// KeyLogWriter: keyLog,
	}

	server := kvserver.Server{
		QuicConfig: quicConf,
		Addr:host,
		TLSConfig: tlsConf,
	}
	server.ListenAndServe()

}

