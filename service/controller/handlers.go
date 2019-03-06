package controller

import (
	"encoding/json"
	"fmt"
	"github.com/charlesparasa/plivotest/plivolibs/auth"
	"net/http"
)

func hello(w http.ResponseWriter,r *http.Request,)  {
  
                w.Write([]byte("Called By API"))
              
}

func signup(w http.ResponseWriter,r *http.Request,)  {
	signUpInput := auth.SignUpParams{}
	if err := json.NewDecoder(r.Body).Decode(&signUpInput); err != nil || signUpInput.Email == "" || signUpInput.Password == "" {
		w.WriteHeader(400)
		w.Write([]byte("Invalid/missing input params"))
		return
	}
	signUpInput.Connection = auth.Realm
	res, resCode, err := Auth0Session.Call(auth.SignupAPIEndpoint, http.MethodPost, auth.ApplicationJSON, signUpInput)
	if err != nil {
			w.Write(res)
		return
	}
	if resCode != http.StatusOK {
		w.Write([]byte("Internal server error please contact the administrator"))
		return
	}
	w.Write(res)
}

func login(w http.ResponseWriter,r *http.Request,)  {
	loginInput := auth.LoginParams{}
	loginInput.ClientID = auth.Auth0C.ClientID
	loginInput.ClientSecret = auth.Auth0C.ClientSecret
	loginInput.Realm = auth.Realm
	loginInput.GrantType = auth.GrantTypePasswordRealm
	loginInput.Scope = auth.ScopeOpenID

	// getting username and password from client
	err := json.NewDecoder(r.Body).Decode(&loginInput)
	if err != nil {
		err = fmt.Errorf("Unbale to unmarshall" , err)
		return
	}

	res, resCode, err := Auth0Session.Call(auth.LoginAPIEndpoint, http.MethodPost, auth.ApplicationJSON, loginInput)
	if err != nil {
		err = fmt.Errorf("No reposne" , err)
		w.Write(res)
	}
	if resCode != http.StatusOK {
		fmt.Println("error" , string(res))
		w.Write(res)
		return
	}
	w.Write(res)
}

var Auth0Session *auth.Auth0

func init() {
	auth.InitAuth0()
	auth.InitJwt()
	Auth0Session = auth.GetAuth0Session()
	if Auth0Session == nil {
		fmt.Print("not able to fetch the session ")
		return
	}
}
