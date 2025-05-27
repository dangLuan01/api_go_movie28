package utilapi

import (
	"encoding/json"
	"net/http"
)
func ResponseWithJson(respone http.ResponseWriter, statusCode int, data interface{}) {
	result, err := json.Marshal(data)
	if err != nil {
		http.Error(respone, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respone.Header().Set("Content-Type", "application/json")
	respone.Header().Set("Access-Control-Allow-Origin", "*")
	respone.WriteHeader(statusCode)
	respone.Write(result)
}