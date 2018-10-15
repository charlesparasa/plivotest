package config

import "context"
//Auth0Conf Auht0 Conf
type Auth0Conf struct {
	Auth0Domain    string `yaml:"domain"`
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	SecretLocation string `yaml:"secret_location"`
	Realm          string `yaml:"realm"`
}

// TContext is the combination of native context and APIContext
type TContext struct {
	context.Context
	APIContext
}

// APIContext contains context of client
type APIContext struct {
	Name string
	Email string
	Phone string
	ClientID       string // ClientID for the clients ID
	Token          string // Token is the api token
	RequestID      string // RequestID - used to track logs across a request-response cycle
	CorrelationID  string // CorrelationID - used to track logs across a user's session
}

