package slickconfig

import (
	"os"
	"strings"
)

type SlickConfiguration struct {
	BaseUrl string `toml:"base-url" comment:"You must supply a base url for slick"`
	JWTPrivateKey []byte `toml:"jwt-private-key" comment:"The private key used to sign the jwt tokens.  Generate using openssl genrsa -out app.rsa 1024"`
	JWTPublicKey []byte `toml:"jwt-public-key" comment:"The public key used to verify the jwt tokens.  Generate using openssl rsa -in app.rsa -pubout > app.rsa.pub"`
}

const (
	locationSystem                       = "/etc/slick.toml"
	locationLocal                        = "./slick.toml"
	ConfigurationEnvironmentVariableName = "SLICKCONF"
)

var (
	Configuration SlickConfiguration
)

func init() {
	Configuration.BaseUrl = "http://localhost:6789"
}

func (c *SlickConfiguration) loadFromFile(location string) {
	
}

func (c *SlickConfiguration) loadFromUrl(location string) {

}

func (c *SlickConfiguration) Load() {
	// load from the location in the environment variable if it's set, next look at the local directory, finally
	// system location
	location, ok := os.LookupEnv(ConfigurationEnvironmentVariableName)
	if !ok {
		if _, err := os.Stat(locationLocal); os.IsNotExist(err) {
			location = locationSystem
		} else {
			location = locationSystem
		}
	}
	c.LoadFromLocation(location)
}

func (c *SlickConfiguration) LoadFromLocation(location string) {
	if strings.HasPrefix(location, "http") {
		c.loadFromUrl(location)
	} else {
		if _, err := os.Stat(location); err == nil {
			c.loadFromFile(location)
		}
	}
}