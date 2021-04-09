package handler

import (
	"MangaLibrary/src/internal/dao"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type WebComponentHandler struct {
	Site         *dao.SitesDAO
	WebComponent *dao.WebComponentDAO
	User         *dao.UsersDAO
	Logger       *zap.Logger
}

func (h WebComponentHandler) NewWebComponent(w http.ResponseWriter, r *http.Request) {
	//userReq := &dto.WebComponent{}
	//err := json.NewDecoder(r.Body).Decode(userReq)
	//if err != nil {
	//	SendError(err.Error(), w)
	//	return
	//}
	//apiKey := r.Header.Get("Authentication-Key")
	//user := &dto.User{ID: userReq.UserID, APIKey: apiKey}
	//err = h.User.UpdateUserAPIKey(user)
	//if err != nil {
	//	SendError(err.Error(), w)
	//	return
	//}
	SendData(nil, "not implemented yet", w)
}

func (h WebComponentHandler) UpdateWebComponent(w http.ResponseWriter, r *http.Request) {
	//userReq := &dto.WebComponent{}
	//err := json.NewDecoder(r.Body).Decode(userReq)
	//if err != nil {
	//	SendError(err.Error(), w)
	//	return
	//}
	//apiKey := r.Header.Get("Authentication-Key")
	//user := &dto.User{ID: userReq.UserID, APIKey: apiKey}
	//err = h.User.UpdateUserAPIKey(user)
	//if err != nil {
	//	SendError(err.Error(), w)
	//	return
	//}
	SendData(nil, "not implemented yet", w)
}

func (h WebComponentHandler) GetWebComponents(w http.ResponseWriter, r *http.Request) {
	siteName := mux.Vars(r)["site_name"]
	site, err := h.Site.GetSite(siteName)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	components, err := h.WebComponent.GetComponentsForSite(site.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}

	SendData(components, "components", w)
}
