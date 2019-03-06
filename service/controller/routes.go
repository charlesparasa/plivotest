package controller

import (
	"github.com/charlesparasa/plivotest/plivolibs/config"
	"github.com/charlesparasa/plivotest/plivolibs/logger"
	"github.com/charlesparasa/plivotest/plivolibs/middleware"
	"net/http"
	"os"
)

func getServicePost() string  {
	if os.Getenv("PORT") == "" {
		return "4030"
	}
	return os.Getenv("PORT")
}

func Start()  {
	middleware.AddNoAuthRoutes(
		"login",
		http.MethodPost,
		"/login",
		login)

	middleware.AddNoAuthRoutes(
		"login",
		http.MethodPost,
		"/signup",
		signup)

		middleware.AddNoAuthRoutes(
			"hello",
			http.MethodGet,
			"hello",
			hello)

	logger.GenericInfo(config.TContext{},"Started Contact service" , logger.FieldsMap{"port":getServicePost()})
	middleware.Start(getServicePost(), "/spike")
}
