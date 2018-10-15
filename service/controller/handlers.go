package controller

import (
	"encoding/json"
	"fmt"
	"github.com/charlesparasa/plivotest/plivolibs/auth"
	"github.com/charlesparasa/plivotest/plivolibs/mysqlManager"
	"github.com/charlesparasa/plivotest/service/model"
	"github.com/gorilla/mux"
	"net/http"
)

func create(w http.ResponseWriter,r *http.Request)  {
	var contact model.Contact
	var res = struct {
		Message string `json:"message"`
	}{
		Message: "Contact Saved",
	}
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		fmt.Println("unbale to decode the request body " , err)
		return
	}
	if contact.Email == "" || contact.Name == ""{
		w.Write([]byte("Please provide email"))
		fmt.Println("missing email/name")
		return
	}
	err = mysqlManager.InsertOne(contact)
	if err != nil {
		 err = fmt.Errorf("Unbale to insert contact" , err)
		 fmt.Println("err", err)
		 res.Message = err.Error()
		 bytes ,err := json.Marshal(res)
		 if err != nil{
			fmt.Println("error", err)
			return
		 }
		 jsonResponse(w, http.StatusBadRequest, string(bytes))
		 return
	}
	bytes ,err := json.Marshal(res)
	if err != nil{
		fmt.Println("error", err)
		jsonResponse(w, http.StatusOK, string(bytes))
	}
}

func getContacts(w http.ResponseWriter,r *http.Request)  {
	var res = struct {
		Message string `json:"message"`
	}{
		Message: "Contact Updated",
	}
	var err error
	fromPage := mux.Vars(r)["from"]
	toPage := mux.Vars(r)["to"]
	contactsInterface, err := mysqlManager.GetData(fromPage, toPage)
	if err != nil {
		fmt.Println("error", err)
		res.Message = err.Error()
		bytes ,err := json.Marshal(res)
		if err != nil{
			fmt.Println("error", err)
			return
		}
		jsonResponse(w, http.StatusBadRequest, string(bytes))
		return
	}
	contactBytes , err := json.Marshal(contactsInterface)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	var contacts []model.Contact
	err = json.Unmarshal(contactBytes ,&contacts)
	if err != nil {
		fmt.Println("error" , err)
		return
	}
	jsonResponse(w, http.StatusOK,string(contactBytes))
}

func deleteContact(w http.ResponseWriter,r *http.Request)  {
	var res = struct {
		Message string `json:"message"`
	}{
		Message: "Contact Deleted",
	}
	email := mux.Vars(r)["email"]
	err := mysqlManager.DeleteContact(email)
	if err != nil {
		fmt.Println("Unable to delete")
		res.Message = err.Error()
		bytes ,err := json.Marshal(res)
		if err != nil{
			fmt.Println("error", err)
			return
		}
		jsonResponse(w, http.StatusBadRequest, string(bytes))
		return
	}
	bytes ,err := json.Marshal(res)
	if err != nil{
		fmt.Println("error", err)
		return
	}
	jsonResponse(w, http.StatusOK, string(bytes))
}

func updateContact(w http.ResponseWriter,r *http.Request)  {
	var res = struct {
		Message string `json:"message"`
	}{
		Message: "Contact Updated",
	}
	var contact model.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil{
		err = fmt.Errorf("Unbale to decode the request body")
		return
	}

	err = mysqlManager.UpdateContact(contact)
	if err != nil {
		fmt.Println("Unable to delete")
		res.Message = err.Error()
		bytes ,err := json.Marshal(res)
		if err != nil{
			fmt.Println("error", err)
			return
		}
		jsonResponse(w, http.StatusBadRequest, string(bytes))
		return
	}
	bytes ,err := json.Marshal(res)
	if err != nil{
		fmt.Println("error", err)
		return
	}
	jsonResponse(w, http.StatusOK, string(bytes))
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
