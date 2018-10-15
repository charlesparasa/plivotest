package middleware

import (
	l "log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	//PlivoAPIToken HTTP Header params that validates user auth token
	PlivoAPIToken = "plivo-api-token"

	//TimeZone of client
	TimeZone = "timezone"

	//Locale of client
	Locale = "locale"

	timeFormat = time.RFC3339

	//RequestID of the parent request
	RequestID = "requestId"

	//CorrelationID of the parent request
	CorrelationID = "correlationId"
)

var routes = make(Routes, 0)

// useMiddleware applies chains of middleware (ie: log, contextWrapper, validateAuth) handler into incoming request
// For example, logging middleware might write the incoming request details to a log
// Note - It applies in reverse order
func useMiddleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// NewRouter provides a mux Router.
// Handles all incoming request who matches registered routes against the request.
func newRouter(subroute string) *mux.Router {

	muxRouter := mux.NewRouter().StrictSlash(true)
	subRouter := muxRouter.PathPrefix(subroute).Subrouter()

	for _, route := range routes {
		subRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return muxRouter
}

//Start - http servers
func Start(port, subroute string) {
	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) // Allowing all origin as of now

	allowedHeaders := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"X-CSRF-Token",
		"X-Auth-Token",
		"Content-Type",
		"processData",
		"contentType",
		"Origin",
		"Authorization",
		"Accept",
		"Client-Security-Token",
		"Accept-Encoding",
		"timezone",
		"locale",
		PlivoAPIToken,
		RequestID,
		CorrelationID})

	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
		"OPTIONS"})

	allowCredential := handlers.AllowCredentials()

	s := &http.Server{
		Addr: ":" + port,
		Handler: handlers.CORS(
			allowedHeaders,
			allowedMethods,
			allowedOrigins,
			allowCredential)(
			context.ClearHandler(
				newRouter(subroute),
			)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	l.Fatal(s.ListenAndServe())
}
