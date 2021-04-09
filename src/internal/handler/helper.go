package handler

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	ERROR_CODE   = 2
	SUCCESS_CODE = 0
)

func SendError(errorMessage string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(Response{
		Message: errorMessage,
		Code:    ERROR_CODE,
	})

}
func SendData(data interface{}, errorMessage string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(Response{
		Message: errorMessage,
		Code:    SUCCESS_CODE,
		Data:    data,
	})
	if err != nil {
		SendError(err.Error(), w)
		return
	}
}

func ServeSocketData() {

}
