package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func SendFile(file string, download bool, w http.ResponseWriter, r *http.Request) {
	if download {
		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(file))
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeFile(w, r, file)
}

func ServeSocketData() {

}
