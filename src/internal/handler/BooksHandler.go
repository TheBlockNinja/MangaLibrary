package handler

import (
	"MangaLibrary/src/internal/dao"
	"net/http"

	"go.uber.org/zap"
)

type BookHandler struct {
	Site   *dao.SitesDAO
	Book   *dao.BooksDAO
	User   *dao.UsersDAO
	Logger *zap.Logger
}

func (h BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authentication-Key")
	u, err := h.User.GetUser(apiKey)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	books, err := h.Book.GetAllUserBooks(u.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(books, "", w)
}
