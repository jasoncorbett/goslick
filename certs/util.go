package certs

import (
	"crypto/tls"
	"crypto/x509"
)

var (
	DemoKeyPair  *tls.Certificate
	DemoCertPool *x509.CertPool
)

func init() {
	var err error
	pair, err := tls.X509KeyPair([]byte(Cert), []byte(Key))
	if err != nil {
		panic(err)
	}
	DemoKeyPair = &pair
	DemoCertPool = x509.NewCertPool()
	ok := DemoCertPool.AppendCertsFromPEM([]byte(Cert))
	if !ok {
		panic("bad certs")
	}
}