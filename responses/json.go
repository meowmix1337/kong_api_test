package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type returnData struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// JSON .
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(returnData{
		Status: statusCode,
		Data:   data,
	})
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// Error .
func Error(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	errJson := json.NewEncoder(w).Encode(returnData{
		Status: statusCode,
		Error:  err.Error(),
	})
	if errJson != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
