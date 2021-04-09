package handler

import (
	"MangaLibrary/src/internal/api"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type LibraryHandler struct {
	Logger *zap.Logger
	//Library *api.Library
	//DB     *sqlx.DB
}

func NewLibraryHandler(logger *zap.Logger) *LibraryHandler {
	return &LibraryHandler{
		Logger: logger,
		//Library: Lib,
	}
}
func (h LibraryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// New Site Data
	case http.MethodGet:
		h.GetLibrary(w, r)
	}
}
func (h LibraryHandler) GetLibrary(w http.ResponseWriter, r *http.Request) {
	dataMap := map[string]string{}
	for k, v := range r.URL.Query() {
		dataMap[k] = v[0]
	}
	siteName := mux.Vars(r)["site_name"]
	Lib := api.NewLibrary()
	Lib.LoadDefaults(h.Logger)
	err := Lib.ProcessSite(siteName, dataMap, nil)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(Lib.Sites[siteName].MetaData, "Finished Getting Data", w)
}

func (h LibraryHandler) CreatePDF(w http.ResponseWriter, r *http.Request) {
	dataMap := map[string]string{}
	for k, v := range r.URL.Query() {
		dataMap[k] = v[0]
	}
	siteName := mux.Vars(r)["site_name"]
	// todo replace this
	Lib := api.NewLibrary()
	Lib.LoadDefaults(h.Logger)
	data, err := Lib.GetPDF(siteName, dataMap)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(data, "finished creating pdfs", w)
}

func (h LibraryHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	fileName, err := url.QueryUnescape(r.URL.Query().Get("file"))
	if err != nil {
		SendError(fmt.Sprintf("failed files %s, does not exist", fileName), w)
		return
	}
	if _, err = os.Stat(fileName); err == nil {
		http.ServeFile(w, r, fileName)
	}
	SendError("file does not exist", w)
}
