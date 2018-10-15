package middleware

import (
	"fmt"
	"github.com/charlesparasa/plivotest/plivolibs/auth"
	"github.com/charlesparasa/plivotest/plivolibs/config"
	"github.com/charlesparasa/plivotest/plivolibs/logger"
	"github.com/dgrijalva/jwt-go"
	jwtM "gopkg.in/square/go-jose.v2/jwt"
)

//ParseJWT Parsing the JWT token for validation
func ParseJWT(jwtString string) (map[string]interface{}, error) {
	ctx := config.TContext{}
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return auth.Data, nil
	})

	if err != nil {
		logger.GenericError(ctx, fmt.Errorf("error parsing the JWT . err-%s", err), nil)
		return nil, err
	}

	if err := token.Claims.Valid(); err != nil {
		logger.GenericError(ctx, fmt.Errorf("error validating the JWT . err-%s", err), nil)
		return nil, err
	}
	return map[string]interface{}(token.Claims.(jwt.MapClaims)), nil
}

// validateToken validates plivo-api-token of incoming request
func validateToken(apiToken string) (string, string, error) {
	var newAccessToken string
	var username string
	var err error
	if err = auth.ValidateAccessJWTSignature(apiToken); err == jwtM.ErrExpired {
		logger.HTTPLog("Token Expired Please contact the administrator")
	}
	if err != nil {
		logger.HTTPLog("Invalid Token Please Check")
		return  "","", err

	}
	return username, newAccessToken, nil
}
