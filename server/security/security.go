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
	certFile := "security/server.crt"
	keyFile := "security/server.key"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// ServerName:         ip,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), config)

	return conn, err
}

func CreateSecureSocket(ip string, port int) (net.Conn, error) {
	certFile := "security/server.crt"
	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", certFile)
		return nil, fmt.Errorf("unable to parse cert")
	}
	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)

	return conn, err
}
