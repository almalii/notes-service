package handler

import (
	"encoding/json"
	"net/http"
)

// TODO: Куда лучше поместить эту функцию?
func WriteJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
