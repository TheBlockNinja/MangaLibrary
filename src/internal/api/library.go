package api

import (
	"fmt"

	"go.uber.org/zap"
)

type Library struct {
	Sites    map[string]*API
	MainPath string
	Progress *Progress
}
type FinishedSites struct {
}

func NewLibrary() *Library {
	return &Library{
		Sites:    map[string]*API{},
		MainPath: "./",
		Progress: &Progress{},
	}
}
func (l *Library) LoadDefaults(logger *zap.Logger) {

}
func (l *Library) GetPDF(name string, search map[string]string) ([]string, error) {
	if _, found := l.Sites[name]; !found {
		return nil, fmt.Errorf("unable to find site")
	}
	currentSite := *l.Sites[name]
	if v, found := search["name"]; found {
		pathName, err := currentSite.GetNameWithURL(fmt.Sprintf("library/%s/%s", name, v))
		if err != nil {
			return nil, err
		}
		return currentSite.CreatePDF(pathName)
	}
	currentURL := currentSite.GetURL(currentSite.SearchURL, search)
	currentSite.Parser.Logger.Info(fmt.Sprintf("Current URL: %s", currentURL))
	pathName, err := currentSite.GetNameWithURL(currentURL)
	if err != nil {
		return nil, err
	}
	return currentSite.CreatePDF(pathName)

}
func (l *Library) AddSite(name string, api *API) {
	l.Sites[name] = api
}
func (l *Library) ProcessSite(name string, search map[string]string, progress *Progress) error {
	if _, found := l.Sites[name]; !found {
		return fmt.Errorf("unable to find site")
	}
	l.Progress = progress
	currentSite := *l.Sites[name]
	currentSite.Parser.Logger.Info(fmt.Sprintf("Current URL: %s", currentSite.SearchURL))
	currentURL := currentSite.GetURL(currentSite.SearchURL, search)
	currentSite.Parser.Logger.Info(fmt.Sprintf("Current URL: %s", currentSite.SearchURL))
	err := currentSite.Process(currentURL, search, l.Progress)
	if err != nil {
		l.Progress.Message = err.Error()
		return err
	}

	return err
}
