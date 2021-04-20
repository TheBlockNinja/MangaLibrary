package handler

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type UserRequest struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}
type UserHandler struct {
	User   *dao.UsersDAO
	Logger *zap.Logger
}

const (
	DefaultMaxJobs = 50
)

func (h UserHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	userReq := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	newUser := &dto.User{
		Name:        userReq.UserName,
		Email:       userReq.Email,
		Password:    userReq.Password,
		APIKey:      "",
		CurrentJobs: 0,
		MaxJobs:     DefaultMaxJobs,
		Age:         userReq.Age,
		IsActive:    true,
		IsAdmin:     false,
	}
	err = h.User.NewUser(newUser)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(newUser, "new api key", w)
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

func (h UserHandler) Signin(w http.ResponseWriter, r *http.Request) {
	userReq := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}

	user, err := h.User.GetUserForName(userReq.UserName, userReq.Password)
	if err != nil {
		SendError("failed to log in", w)
		return
	}
	SendData(user, "logged in", w)
}
