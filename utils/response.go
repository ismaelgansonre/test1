package utils

import (
	"encoding/json"
	"net/http"
)

// SendJSONResponse sends a JSON response with the given status code
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
