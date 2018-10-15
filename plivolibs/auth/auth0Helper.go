package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charlesparasa/plivotest/plivolibs/config"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	clientIDConst     = "clientid"
	clientSecretConst = "clientsecret"
)

// TokenRefresher - Use this endpoint to refresh an Access Token using the Refresh Token you got during authorization.
type TokenRefresher struct {
	GrantType    string `json:"grant_type"`    // REQUIRED - Denotes the flow you are using. For Client Credentials use
	ClientID     string `json:"client_id"`     // REQUIRED - Your application's Client ID.
	ClientSecret string `json:"client_secret"` // REQUIRED - Your application's Client Secret. Required when the Token Endpoint Authentication Method  field at your Client Settings is Post or Basic
	RefreshToken string `json:"refresh_token"` // REQUIRED - The Refresh Token to use.
}

// Auth0client auth0 client
type Auth0client struct {
	ClientID     string
	ClientSecret string
}
type Output struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

//OauthLoginResponse structure to capture login response.
type OauthLoginResponse struct {
	IDToken      string `json:"id_token,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}

//Auth0C Varible which the Auth0 conf
var Auth0C Auth0Conf

//Auth0Conf Auht0 Conf
type Auth0Conf struct {
	Auth0Domain    string `yaml:"domain"`
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	SecretLocation string `yaml:"secret_location"`
	Realm          string `yaml:"realm"`
}


var A *Auth0

func InitAuth0()*Auth0{
	auth0Conf, err := config.GetAuth0Config()
	if err != nil {
		err := fmt.Errorf("GetAuth0Config: error getting the auth0 conf: %v", err)
		fmt.Println("error" , err)
	}
	Auth0C.SecretLocation = auth0Conf.SecretLocation
	Auth0C.Realm = auth0Conf.Realm
	Auth0C.ClientID = auth0Conf.ClientID
	Auth0C.ClientSecret = auth0Conf.ClientSecret
	Auth0C.Auth0Domain = auth0Conf.Auth0Domain

	A = NewAuth0(auth0Conf.Auth0Domain, "")
	return A
}

func GetAuth0Session() *Auth0 {
	return A
}

// SignUpParams struct
type SignUpParams struct {
	ClientID   string `json:"client_id"`  // REQUIRED
	Email      string `json:"email"`      // REQUIRED
	Password   string `json:"password"`   // REQUIRED
	Connection string `json:"connection"` // REQUIRED
	Username   string `json:"username"`   // OPTIONAL (Required if we have set it as required fields in auth0 connections)
}

//SignUpResponseError - standard authentication auth0 error response
type SignUpResponseError struct {
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description interface{} `json:"description"`
	StatusCode  int         `json:"statusCode"`
	Error       string      `json:"error"`
}

//SignUpResponse - standard response of Auth0 signup API
type SignUpResponse struct {
	ID            string `json:"_id"`
	EmailVerified bool   `json:"email_verified"`
	Email         string `json:"email"`
	Username      string `json:"username"`
}

//LoginParams Auth0 standard login struct
type LoginParams struct {
	//REQUIRED
	//The ClientID of your client
	ClientID string `json:"client_id"`

	//REQUIRED
	//The Client Secret of your client
	ClientSecret string `json:"client_secret"`

	//REQUIRED
	//email of the user to login
	Username string `json:"username"`

	//REQUIRED
	//Password of the user to login
	Password string `json:"password"`

	//REQUIRED
	//The name of the realm/connection to use for login
	Realm string `json:"realm"`

	//REQUIRED
	//Set to 'http://auth0.com/oauth/grant-type/password-realm' to authenticate using username/password with connection
	//OR urn:ietf:params:oauth:grant-type:jwt-bearer to authenticate using an id_token
	GrantType string `json:"grant_type"`

	//Set to openid to retrieve also an id_token,
	//leave null to get only an access_token
	Scope string `json:"scope"`
}

//LoginResponse structure to capture login response.
type LoginResponse struct {
	IDToken      string `json:"id_token,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}

//User - standard Auth0 user https://auth0.com/docs/api/management/v2#!/Users/get_users
type User struct {
	Connection    string   `json:"connection"`         //REQUIRED
	Email         string   `json:"email,omitempty"`    //REQUIRED
	Username      string   `json:"username,omitempty"` //REQUIRED
	Password      string   `json:"password,omitempty"` //REQUIRED
	PhoneNumber   string   `json:"phone_number,omitempty"`
	Blocked       string   `json:"blocked,omitempty"`
	UserMetadata  struct{} `json:"user_metadata,omitempty"`
	Name          string   `json:"name,omitempty"`
	EmailVerified bool     `json:"email_verified,omitempty"`
	VerifyEmail   bool     `json:"verify_email,omitempty"`
	PhoneVerified bool     `json:"phone_verified,omitempty"`
	AppMetadata   struct{} `json:"app_metadata,omitempty"`
	RefreshToken  string   `json:"refreshToken,omitempty"`
}

//NewAuth0 returns a new instance of auth0
func NewAuth0(apiBasePath, token string) *Auth0 {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	return &Auth0{
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: netTransport,
		},
		apiBasePath: apiBasePath,
		Token:       token,
	}
}

//Call makes a https call to auth0
func (auth0Session Auth0) Call(apiEndPoint, method, contentType string, body interface{}) ([]byte, int, error) {
	var encodedBody []byte
	urlStr := auth0Session.apiBasePath + apiEndPoint
	if _, ok := body.([]byte); !ok {
		encodedBody, _ = json.Marshal(body)
	} else {
		encodedBody = body.([]byte)
	}
	fmt.Println("url ", urlStr)
	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(encodedBody))
	req.Header.Add(Authorization, AuthorizationType+auth0Session.Token)
	req.Header.Add(ContentTypeHeader, contentType)
	resp, err := auth0Session.client.Do(req)
	if err != nil {
		return nil, http.StatusMethodNotAllowed, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusMethodNotAllowed, err
	}
	return resBody, resp.StatusCode, nil
}

// SecretProvider will provide everything
// needed retrieve the secret.
type SecretProvider interface {
	GetSecret(r *http.Request) (interface{}, error)
}

// SecretProviderFunc simple wrappers to provide
// secret with functions.
type SecretProviderFunc func(*http.Request) (interface{}, error)

// GetSecret implements the SecretProvider interface.
func (f SecretProviderFunc) GetSecret(r *http.Request) (interface{}, error) {
	return f(r)
}

// NewKeyProvider provide a simple passphrase key provider.
func NewKeyProvider(key interface{}) SecretProvider {
	return SecretProviderFunc(func(_ *http.Request) (interface{}, error) {
		return key, nil
	})
}

// Configuration contains
// all the information about the
// Auth0 service.
type Configuration struct {
	secretProvider SecretProvider
	expectedClaims jwt.Expected
	signIn         jose.SignatureAlgorithm
}