package handler

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type SiteRequest struct {
	Site          *dto.Site           `json:"site"`
	WebComponents []*dto.WebComponent `json:"web_components"`
}
type SiteHandler struct {
	Site         *dao.SitesDAO
	WebComponent *dao.WebComponentDAO
	User         *dao.UsersDAO
	Logger       *zap.Logger
}

func (h SiteHandler) NewSite(w http.ResponseWriter, r *http.Request) {
	siteRequest := &SiteRequest{}
	err := json.NewDecoder(r.Body).Decode(siteRequest)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	err = h.Site.NewSite(siteRequest.Site)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	for _, wc := range siteRequest.WebComponents {
		wc.SiteID = siteRequest.Site.ID
		err = h.WebComponent.NewComponent(wc)
		if err != nil {
			SendError(err.Error(), w)
			return
		}
	}
	SendData(siteRequest, "new site data", w)
}

func (h SiteHandler) GetSite(w http.ResponseWriter, r *http.Request) {
	siteName := mux.Vars(r)["site_name"]
	site, err := h.Site.GetSite(siteName)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(site, "site data", w)

}

func (h SiteHandler) GetAllSites(w http.ResponseWriter, r *http.Request) {
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	site, err := h.Site.GetAllSites(user.Age)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(site, "all site data", w)

}
