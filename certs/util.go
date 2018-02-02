package certs

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/serussell/logxi/v1"
)

var (
	DemoKeyPair  *tls.Certificate
	DemoCertPool *x509.CertPool
	logger log.Logger
)

func init() {
	var err error
	logger = log.New("certs")
	pair, err := tls.X509KeyPair([]byte(Cert), []byte(Key))
	if err != nil {
		logger.Fatal("Unable to create cert / key pair", err)
	}
	DemoKeyPair = &pair
	DemoCertPool = x509.NewCertPool()
	ok := DemoCertPool.AppendCertsFromPEM([]byte(Cert))
	if !ok {
		logger.Fatal("The certs provided are bad in some way, unable to add to cert pool.")
	}
}

