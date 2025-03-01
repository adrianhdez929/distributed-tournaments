package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

func CreateSecureSocketListener(port int) (net.Listener, error) {
	certFile := "server.crt"
	keyFile := "server.key"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	conn, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), config)

	return conn, err
}

func CreateSecureSocket(ip string, port int) (net.Conn, error) {
	certFile := "server.crt"
	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", certFile)
	}
	config := &tls.Config{RootCAs: certPool}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)

	return conn, err
}
