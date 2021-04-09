package dao

import (
	"MangaLibrary/src/internal/dto"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SitesDAO struct {
	DB *sqlx.DB
}

func (s *SitesDAO) NewSite(site *dto.Site) error {
	if site.BaseURL == "" || site.Name == "" || site.SearchURL == "" {
		return fmt.Errorf("missing base_url, name, or search_url")
	}
	if site.BasePath == "" {
		site.BasePath = "library/" + strings.ReplaceAll(site.Name, " ", "")
	}
	insertStmt := "insert into manga_library.sites(name,base_url,search_url,base_path,use_sub_path,meta_data,min_age) VALUES(?,?,?,?,?,?,?)"
	results, err := s.DB.Exec(insertStmt, site.Name, site.BaseURL, site.SearchURL, site.BasePath, site.UseSubPath, site.MetaData, site.MinAge)
	if err != nil {
		return err
	}
	id, _ := results.LastInsertId()
	site.ID = int(id)
	return nil
}

func (s *SitesDAO) GetSite(name string) (*dto.Site, error) {
	site := &dto.Site{}
	err := s.DB.Get(site, "select * from manga_library.sites where name = ?", name)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (s *SitesDAO) GetAllSites(age int) ([]*dto.Site, error) {
	var siteList []*dto.Site
	rows, err := s.DB.Queryx("select * from manga_library.sites where min_age <= ? order by name", age)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		site := &dto.Site{}
		err = rows.StructScan(site)
		if err != nil {
			return nil, err
		}
		siteList = append(siteList, site)
	}
	return siteList, nil
}

func (s *SitesDAO) UpdateSite() {

}
func (s *SitesDAO) DeleteSite() {

}
