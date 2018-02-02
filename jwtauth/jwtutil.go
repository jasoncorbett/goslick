package jwtauth

import (
	"github.com/SermoDigital/jose/jws"
	"crypto/rsa"
	"github.com/SermoDigital/jose/crypto"
	"github.com/jasoncorbett/goslick/certs"
	"github.com/serussell/logxi/v1"
	"fmt"
	"errors"
	"google.golang.org/grpc/credentials"
	"context"
)

var (
	JwtRSAPrivateKey *rsa.PrivateKey
	JwtRSAPublicKey *rsa.PublicKey
	logger log.Logger
)

func init() {
	logger = log.New("jwtauth")
	var err error
	JwtRSAPrivateKey, err = crypto.ParseRSAPrivateKeyFromPEM([]byte(certs.JwtKey))
	if err != nil {
		logger.Fatal("Error creating Private RSA Key for Jwt Auth Signing", err)
	}
	JwtRSAPublicKey, err = crypto.ParseRSAPublicKeyFromPEM([]byte(certs.JwtPublicKey))
	if err != nil {
		logger.Fatal("Error creating Public RSA Key for Jwt Auth Signing", err)
	}
}

type jwt struct {
	token string
}

func NewCredential() (credentials.PerRPCCredentials, error) {
	token, _ := CreateJWT()
	return jwt{token}, nil
}

func (j jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", j.token),
	}, nil
}

func (j jwt) RequireTransportSecurity() bool {
	return true
}


func CreateJWT() (string, error) {
	claims := jws.Claims{}
	permissions := [...]string{"hello", "world"}
	claims.Set("p", permissions)
	signMethod := jws.GetSigningMethod("RS256")
	token := jws.NewJWT(claims, signMethod)
	byteToken, err := token.Serialize(JwtRSAPrivateKey)
	if err != nil {
		logger.Fatal("Error signing the jwt. ", err)
		return  "", err
	}

	return string(byteToken), nil
}

func PermissionsFromJWT(rawtoken string) ([]string, error) {
	token, err := jws.ParseJWT([]byte(rawtoken))

	if err != nil {
		logger.Error("Error parsing jwt.txt: ", err)
		return make([]string, 0), err
	}

	if err = token.Validate(JwtRSAPublicKey, crypto.SigningMethodRS256); err != nil {
		logger.Error("JWT Validation Failed: ", err)
		return make([]string, 0), err
	}

	uncasted := token.Claims().Get("p")
	uncastedArray, ok := uncasted.([]interface{})
	if ! ok {
		errorText := fmt.Sprintf("Permissions of wrong type: %#v", uncasted)
		logger.Error(errorText)
		return make([]string, 0), errors.New(errorText)
	}

	permissions := make([]string, len(uncastedArray))
	for i, perm := range uncastedArray {
		permissions[i] = perm.(string)
	}
	return permissions, nil
}
