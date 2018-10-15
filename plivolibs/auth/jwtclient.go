package auth

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"io/ioutil"
	"net/http"
	"os"
)

var secretProvider auth0.SecretProvider

//Data for return the secret
var Data []byte

// AccessTokenExtractorFunc function conforming
// to the RequestTokenExtractor interface.
type accessTokenExtractor string

//GetPemFileLocation get the location of the pem file in the container or local machine from the DB
func GetPemFileLocation() (string, error) {
	loc := Auth0C.SecretLocation
	return loc, nil
}

func InitJwt() {
	//Creates a configuration with the Auth0 information
	loc, err := GetPemFileLocation()
	if err != nil {
		fmt.Errorf("Unable to fetch the secret location", err)
		return
	}

	data, err := ioutil.ReadFile(loc)
	if err != nil {
		err = fmt.Errorf("impossible to read key form disk pem file, hence exiting application gracefully")
		fmt.Println("Error reading the pem file from the disk ignoring")
		os.Exit(0)
	}
	Data = data
	secret, err := loadPublicKey(data)
	if err != nil {
		fmt.Println("invalid provided key, hence exiting application gracefully ignoring", err)
		os.Exit(0)
	}
	secretProvider = auth0.NewKeyProvider(secret)

}

// Extract calls f(r)
func (t accessTokenExtractor) Extract(_ *http.Request) (*jwt.JSONWebToken, error) {
	return jwt.ParseSigned(string(t))
}

// loadPublicKey loads a public key from PEM/DER-encoded data.
func loadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}

// ValidateAccessJWTSignature validates JWT signature of a client token using algorithm RS256
func ValidateAccessJWTSignature(accessToken string) error {
	validator := getValidator(accessToken)
	_, err := validator.ValidateRequest(nil)
	return err
}

// getValidator constructs a jwt validate
func getValidator(accessToken string) *auth0.JWTValidator {
	configuration := auth0.NewConfiguration(secretProvider, []string{Auth0C.ClientID}, Auth0C.Auth0Domain+"/", jose.RS256)
	return auth0.NewValidator(configuration, accessTokenExtractor(accessToken))
}
