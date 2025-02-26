package api

import (
	"WasaText/service/consts"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type APIError struct {
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

func NewAPIError(message string, statusCode int) *APIError {
	return &APIError{Message: message, StatusCode: statusCode}
}

func HandleError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if ok := errors.As(err, &apiErr); ok {
		w.WriteHeader(apiErr.StatusCode)
		http.Error(w, apiErr.Message, apiErr.StatusCode)
		log.Printf("%sError: %s%s", consts.Red, apiErr.Error(), consts.Reset)
	} else {

		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Unexpected Error: %s", err.Error())
	}
}
