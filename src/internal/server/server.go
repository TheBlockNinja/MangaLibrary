package server

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/handler"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func StartServer(db *sqlx.DB, logger *zap.Logger) error {
	logger.Info("starting server on port 8080...")
	router := mux.NewRouter()

	LibraryHandler := handler.NewLibraryHandler(logger)

	userDAO := &dao.UsersDAO{DB: db}
	siteDao := &dao.SitesDAO{DB: db}
	webComponentDao := &dao.WebComponentDAO{DB: db}
	jobDao := &dao.JobDAO{DB: db}
	bookDao := &dao.BooksDAO{DB: db}
	UserHandler := handler.UserHandler{User: userDAO, Logger: logger}
	SiteHandler := handler.SiteHandler{
		Site:         siteDao,
		WebComponent: webComponentDao,
		User:         userDAO,
		Logger:       logger,
	}
	WebComponentHandler := handler.WebComponentHandler{
		Site:         siteDao,
		WebComponent: webComponentDao,
		User:         userDAO,
		Logger:       logger,
	}

	JobHandler := handler.JobHandler{
		Site:   siteDao,
		Job:    jobDao,
		User:   userDAO,
		Logger: logger,
	}

	BookHandler := handler.BookHandler{
		Site:   siteDao,
		Book:   bookDao,
		User:   userDAO,
		Logger: logger,
	}
	mid := MiddleWare{users: userDAO, logger: logger}
	//router.Handle("/site_data", jobHandler)
	testHandler := handler.SocketHandler{Logger: logger}
	router.HandleFunc("/socket/library/{site_name}", testHandler.Handle)
	router.Handle("/library/{site_name}", LibraryHandler)
	router.HandleFunc("/pdf/{site_name}", LibraryHandler.CreatePDF) // groupby="chapter/volume"
	router.HandleFunc("/file", LibraryHandler.GetFile)

	//v2
	/*
		router.Handle("v2/jobs").Method("Get","POST")
		router.Handle("v2/library/{library_type}").Method("Get","delete")

		router.Handle("v2/sites").Method("Get","put","post","delete)
		router.Handle("v2/sites/list").Method("Get","put","post","delete)
		router.Handle("v2/sites").Method("Get","put","post","delete)
		router.Handle("v2/sites/test").Method("Get","put","post")
		router.Handle("v2/files").Method("Get","put","post")
	*/
	router.HandleFunc("/v2/user/api_key", UserHandler.NewAPIKey)

	router.HandleFunc("/v2/site", SiteHandler.NewSite).Methods(http.MethodPost)
	router.HandleFunc("/v2/site", SiteHandler.GetAllSites).Methods(http.MethodGet)
	router.HandleFunc("/v2/site/{site_name}", SiteHandler.GetSite).Methods(http.MethodGet)

	router.HandleFunc("/v2/component/{site_name}", WebComponentHandler.GetWebComponents).Methods(http.MethodGet)

	router.HandleFunc("/v2/job", JobHandler.NewJob).Methods(http.MethodPost)
	router.HandleFunc("/v2/job", JobHandler.GetJobs).Methods(http.MethodGet)

	router.HandleFunc("/v2/books/all", BookHandler.GetAllBooks).Methods(http.MethodGet)

	router.Use(CORS, mid.Auth)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Error("Failed to start server", zap.Error(err))
		return err
	}
	return nil
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Next
		next.ServeHTTP(w, r)
		return
	})
}

type MiddleWare struct {
	logger *zap.Logger
	users  *dao.UsersDAO
}

func (mid *MiddleWare) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Next
		apiKey := r.Header.Get("Authentication-Key")
		u, err := mid.users.GetUser(apiKey)
		if err != nil {
			mid.logger.Error("failed getting user", zap.Error(err))
			handler.SendError("invalid Authentication-Key", w)
			return
		}
		mid.logger.Info(fmt.Sprintf("%s : accessing %s", u.Name, r.URL.String()))
		next.ServeHTTP(w, r)
		return
	})
}
