package controller

import "net/http"

// response: an aribtary function to make response
func jsonResponse(w http.ResponseWriter, statusCode int, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(data))
	return
}
