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
	middleware.AddRoute(
		"Create Contact",
		http.MethodPost,
		"/create",
		create)

	middleware.AddRoute(
		"Get Contact",
		http.MethodGet,
		"/getContacts/{from}/{to}",
		getContacts)

	middleware.AddRoute(
		"get Contact By Email",
		http.MethodGet,
		"/contact/{email}/email",
		getContactByEmail)

	middleware.AddRoute(
		"get Contact By Name",
		http.MethodGet,
		"/contacts/{name}/name",
		getContactByName)


	middleware.AddRoute(
		"Get Contact",
		http.MethodDelete,
		"/delete/{email}",
		deleteContact)

	middleware.AddRoute(
		"Get Contact",
		http.MethodPatch,
		"/update",
		updateContact)

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

	logger.GenericInfo(config.TContext{},"Started Contact service" , logger.FieldsMap{"port":getServicePost()})
	middleware.Start(getServicePost(), "/contacts")
}