package handler

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type UserRequest struct {
	UserID int `json:"user_id"`
}
type UserHandler struct {
	User   *dao.UsersDAO
	Logger *zap.Logger
}

func (h UserHandler) NewAPIKey(w http.ResponseWriter, r *http.Request) {
	userReq := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	apiKey := r.Header.Get("Authentication-Key")
	user := &dto.User{ID: userReq.UserID, APIKey: apiKey}
	err = h.User.UpdateUserAPIKey(user)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(user.APIKey, "new api key", w)
}
