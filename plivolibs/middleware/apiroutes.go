package middleware

import (
	"fmt"
	"github.com/twinj/uuid"
	"net/http"
)

//AddNoAuthRoutes - Route without any Auth
func AddNoAuthRoutes(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, logRequest)}
	routes = append(routes, r)

}

// AddRoute is to create routes with ACL enforcer
func AddRoute(methodName, methodType, mRoute string,handlerFunc http.HandlerFunc) {
	r := route{
		Name:                   methodName,
		Method:                 methodType,
		Pattern:                mRoute,
		HandlerFunc:            useMiddleware(handlerFunc, logRequest, validateContext, createContext),
	}
	routes = append(routes, r)
}

// validateContext incoming request context
func validateContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			err            error
			newAccessToken string
		)
		ctx := r.Context()
		apiCtx, ok := ctx.Value(APICtx).(APIContext)
		tempContext := TContext{APIContext: apiCtx}
		if !ok {
			err = fmt.Errorf("Invalid context", tempContext)
			return
		}

		if len(apiCtx.Token) == 0 {
			err = fmt.Errorf("Invalid Token")
			fmt.Println("Error", err)
			return
		}

		apiCtx.Name, newAccessToken, err = validateToken(apiCtx.Token)
		if err != nil {
			err = fmt.Errorf("Invalid Token" , err)
			fmt.Println("Error", err)
			w.WriteHeader(400)
			w.Write([]byte("invalid token"))
			return
		}
		if len(newAccessToken) > 0 {
			w.Header().Set(PlivoAPIToken, newAccessToken)
		}
		ctx = WithAPICtx(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// createContext generates plivo API context for incoming request
// and appends in request Context
func createContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		ctx := r.Context()
		requestID := header.Get(RequestID)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}
		correlationID := header.Get(CorrelationID)
		if correlationID == "" {
			correlationID = uuid.NewV4().String()
		}

		token, client := header.Get(PlivoAPIToken), header.Get("clientID")
		apiCtx := APIContext{
			ClientID:      client,
			Token:         token,
			RequestID:     requestID,
			CorrelationID: correlationID,
		}
		ctx = WithAPICtx(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
