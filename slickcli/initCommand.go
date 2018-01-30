package main

import (
	"github.com/codegangsta/cli"
	"github.com/serussell/logxi/v1"
	"github.com/jasoncorbett/goslick/slickconfig"
	"crypto/rsa"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
	"fmt"
)

var (
	InitCommand = cli.Command{
		Name: "init",
		Flags: []cli.Flag {
			cli.StringFlag{
				Name: "base-url",
				Value: slickconfig.Configuration.BaseUrl,
			},
		},
		Action: InitializeConfiguration,
	}
)

func InitializeConfiguration(c *cli.Context) {
	logger := log.New("slickcli.init")
	if c.IsSet("base-url") {
		slickconfig.Configuration.BaseUrl = c.String("base-url")
	}
	logger.Debug("Generating key...")
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		logger.Fatal("Unable to generate cryptographic key", err)
		return
	}
	block := &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	slickconfig.Configuration.JWTPrivateKey = string(pem.EncodeToMemory(block))
	public := key.Public()
	buffer, err := x509.MarshalPKIXPublicKey(public)
	block = &pem.Block{
		Type: "RSA PUBLIC KEY",
		Bytes: buffer,
	}
	slickconfig.Configuration.JWTPublicKey = string(pem.EncodeToMemory(block))
	configContent, err := slickconfig.Configuration.ToBytes()
	if err != nil {
		logger.Fatal("Problem generating configuration:", err)
	}
	fmt.Println(string(configContent))
}
