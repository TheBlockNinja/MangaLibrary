package handler

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type LibraryHandlers struct {
	User    *dao.UsersDAO
	Library *dao.LibraryDAO
	Logger  *zap.Logger
}

func (h LibraryHandlers) NewItem(w http.ResponseWriter, r *http.Request) {
	libraryReq := &dto.Library{}
	err := json.NewDecoder(r.Body).Decode(libraryReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	libraryReq.UserID = user.ID

	err = h.Library.NewElement(libraryReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(libraryReq, "created library element", w)
}

func (h LibraryHandlers) GetCollections(w http.ResponseWriter, r *http.Request) {
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	collections, err := h.Library.GetCollections(user.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(collections, "", w)
}

func (h LibraryHandlers) GetCollectionItems(w http.ResponseWriter, r *http.Request) {
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		SendError("no collection provided", w)
		return
	}
	collectionItems, err := h.Library.GetItemsForCollection(collection, user.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(collectionItems, "collection items", w)
}

func (h LibraryHandlers) UpdateCollectionItem(w http.ResponseWriter, r *http.Request) {
	libraryReq := &dto.Library{}
	err := json.NewDecoder(r.Body).Decode(libraryReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	libraryReq.UserID = user.ID

	err = h.Library.UpdateItem(libraryReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(libraryReq, "updated library element", w)
}
func (h LibraryHandlers) DeleteCollectionItem(w http.ResponseWriter, r *http.Request) {
	libraryReq := &dto.Library{}
	err := json.NewDecoder(r.Body).Decode(libraryReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	libraryReq.UserID = user.ID

	err = h.Library.DeleteItem(libraryReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(libraryReq, "updated library element", w)
}

func (h LibraryHandlers) GetLibraryBook(w http.ResponseWriter, r *http.Request) {
	libraryReq := &dto.Library{}
	err := json.NewDecoder(r.Body).Decode(libraryReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	libraryReq.UserID = user.ID

	book, err := h.Library.GetLibraryItem(user.ID, libraryReq.BookID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(book, "get library element", w)
}
