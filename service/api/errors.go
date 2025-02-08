package api

import (
	"WasaText/service/consts"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// APIError is a custom error structure
type APIError struct {
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

// Helper to create APIError
func NewAPIError(message string, statusCode int) *APIError {
	return &APIError{Message: message, StatusCode: statusCode}
}

func HandleError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if ok := errors.As(err, &apiErr); ok {
		fmt.Printf("Error: %s\n", apiErr.StatusCode)
		// Custom API error
		w.WriteHeader(apiErr.StatusCode)
		http.Error(w, apiErr.Message, apiErr.StatusCode)
		log.Printf("%sError: %s%s", consts.Red, apiErr.Error(), consts.Reset)
	} else {
		// Internal server error for unexpected errors
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Unexpected Error: %s", err.Error())
	}
}
