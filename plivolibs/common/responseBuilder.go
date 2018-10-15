package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func renderJSON(w http.ResponseWriter, status int, val interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(val)
	w.Write(b)
}

//WriteJSONTapResponseTo generate the api response for Tap
func WriteJSONResponse(w http.ResponseWriter, data interface{}, code *string, message *string) {
	if code == nil {
		codeStr := strconv.Itoa(http.StatusOK)
		code = &codeStr
	}
	if message == nil {
		messageStr := http.StatusText(http.StatusOK)
		message = &messageStr
	}
	renderJSON(w, http.StatusOK, StandardResponse{
		Meta: Metadata{
			Code:    *code,
			Message: *message,
		},
		Data: data,
	})
}

//WriteJSONErrorTapResponse to Generate the error response for tap api
func WriteJSONErrorResponse(w http.ResponseWriter, httpStatus int, err error, data interface{}, code *string, message *string) {
	if code == nil {
		codeStr := strconv.Itoa(httpStatus)
		code = &codeStr
	}
	if message == nil {
		messageStr := http.StatusText(httpStatus)
		message = &messageStr
	}
	renderJSON(w, httpStatus, StandardResponse{
		Meta: Metadata{
			Code:    *code,
			Message: *message,
		},
		Data: data,
	})
}
