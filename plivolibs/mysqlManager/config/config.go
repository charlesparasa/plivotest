package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type appConfig struct {
	Mysql struct {
		Host         string `json:"host"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		DatabaseName string `json:"databaseName"`
	} `json:"mysql"`
}

// AppConfig for application config
var AppConfig appConfig

//Init is for config init
func Init() {
	configStr := os.Getenv("automation_config")
	if len(configStr) == 0 {
		err := errors.New("config.init: Please set ENV variable")
		panic(err)
	}
	decErr := json.NewDecoder(strings.NewReader(configStr)).Decode(&AppConfig)
	if decErr != nil {
		err := fmt.Errorf("config.init: error decoding env variable value %s; err: %v", configStr, decErr)
		panic(err)
	}
}
