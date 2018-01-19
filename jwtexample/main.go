package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "crypto/rsa"
)

import "github.com/SermoDigital/jose/jws"
import (
	"github.com/SermoDigital/jose/crypto"
	"strings"
)

const (
    privKeyPath = "keys/app.rsa"     // openssl genrsa -out keys/app.rsa 1024
    pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in keys/app.rsa -pubout > keys/app.rsa.pub
)

var signKey []byte
var privKey *rsa.PrivateKey
var pubKey *rsa.PublicKey

func init() {
    var err error
    signKey, err = ioutil.ReadFile(privKeyPath)
    if err != nil {
        log.Fatal("Error reading private key")
        os.Exit(1)
    }
    privKey, err = crypto.ParseRSAPrivateKeyFromPEM(signKey)
    if err != nil {
		log.Fatal("Error parsing private key.  ", err)
		os.Exit(1)
    }

	signKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		os.Exit(1)
	}
	pubKey, err = crypto.ParseRSAPublicKeyFromPEM(signKey)
	if err != nil {
		log.Fatal("Error parsing public key.  ", err)
		os.Exit(1)
	}
}

func main() {
	filecontents, err := ioutil.ReadFile("jwt.txt")
	if err != nil {
		log.Fatal("Error reading jwt.txt.  ", err)
		os.Exit(1)
	}
    token, err := jws.ParseJWT(filecontents)
    if err != nil {
    	log.Fatal("Error parsing jwt.txt: ", err)
    	os.Exit(1)
	}
	if err = token.Validate(pubKey, crypto.SigningMethodRS256); err != nil {
		log.Fatal("JWT Validation Failed: ", err)
		os.Exit(1)
	}
	fmt.Println("JWT Validated!")
	uncasted := token.Claims().Get("p")
	uncastedArray, ok := uncasted.([]interface{})
	if ! ok {
		log.Fatalf("Permissions of wrong type: %#v", uncasted)
		os.Exit(1)
	}
	permissions := make([]string, len(uncastedArray))
	for i, perm := range uncastedArray {
		permissions[i] = perm.(string)
	}
	fmt.Printf("Permissions Found: \"%s\"\n", strings.Join(permissions, "\", \""))
}

func createJWT() string {
    claims := jws.Claims{}
    permissions := [...]string{"hello", "world"}
    claims.Set("p", permissions)
    signMethod := jws.GetSigningMethod("RS256")
    token := jws.NewJWT(claims, signMethod)
    byteToken, err := token.Serialize(privKey)
    if err != nil {
        log.Fatal("Error signing the key. ", err)
        os.Exit(1)
    }

    return string(byteToken)
}

