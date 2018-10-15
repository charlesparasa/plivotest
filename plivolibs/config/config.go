package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

//Auth0 https client wrapper for all auth0 API call
type Auth0 struct {
	client      *http.Client
	apiBasePath string
	Token       string
}

var Auth0C *Auth0Conf
func GetAuth0ConfigFromYaml() (*Auth0Conf, error)  {
	var conf *Auth0Conf
	c, err :=ioutil.ReadFile("./plivolibs/config/auth0conf.yaml")
	if err != nil {
		fmt.Println("GetConnection", err)
		err = fmt.Errorf("unable to read the conf ", err)
		return nil, err
	}

	err = yaml.Unmarshal(c, &conf)
	if err != nil {
		err = fmt.Errorf("unable to Unmarshall the conf ", err)
		return nil, err
	}
	Auth0C = conf
	return conf, nil
}

//GetAuth0Config Get the Auht0Conf from DB
func GetAuth0Config() (*Auth0Conf, error) {
	return GetAuth0ConfigFromYaml()
}

