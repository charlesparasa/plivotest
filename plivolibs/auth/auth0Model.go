package auth

import "net/http"

const (

	// GrantTypePasswordRealm ...
	GrantTypePasswordRealm = "http://auth0.com/oauth/grant-type/password-realm"

	// ScopeOpenID ....
	ScopeOpenID = "openid offline_access"

	// Realm - Database connection
	Realm = "Username-Password-Authentication"

	// LoginAPIEndpoint login api end point
	LoginAPIEndpoint = "/oauth/token"

	//SignupAPIEndpoint signup api endpoint
	SignupAPIEndpoint = "/dbconnections/signup"

	// UserAPIEndpoint user related api endpoint
	UserAPIEndpoint = "/api/v2/users"

	//AuthorizationType - standard authorization header for all auth0 API
	AuthorizationType = "Bearer "

	//ApplicationJSON application/json header type
	ApplicationJSON = "application/json"

	//ContentTypeHeader a header type
	ContentTypeHeader = "Content-Type"

	//Authorization constant string
	Authorization = "Authorization"
)

//Auth0 https client wrapper for all auth0 API call
type Auth0 struct {
	client      *http.Client
	apiBasePath string
	Token       string
}


