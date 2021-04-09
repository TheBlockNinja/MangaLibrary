package dao

import (
	"MangaLibrary/src/internal/dto"

	"github.com/jmoiron/sqlx"
)

type WebComponentDAO struct {
	DB *sqlx.DB
}

func (s *WebComponentDAO) NewComponent(webComponent *dto.WebComponent) error {
	insertStmt := "insert into manga_library.web_component(site_id,name,tag,attribute,value,is_link,is_download,link_attributes,element_data,parent,delay,meta_data,reverse) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	results, err := s.DB.Exec(insertStmt, webComponent.SiteID, webComponent.Name, webComponent.Tag, webComponent.Attribute, webComponent.Value, webComponent.IsLink, webComponent.IsDownload, webComponent.LinkAttributes, webComponent.ElementData, webComponent.Parent, webComponent.Delay, webComponent.MetaData, webComponent.Reverse)
	if err != nil {
		return err
	}
	id, _ := results.LastInsertId()
	webComponent.ID = int(id)
	return nil
}

func (s *WebComponentDAO) GetComponentsForSite(siteId int) ([]*dto.WebComponent, error) {
	var webComponentList []*dto.WebComponent
	rows, err := s.DB.Queryx("select * from manga_library.web_component where site_id = ? order by parent,id", siteId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		webComponent := &dto.WebComponent{}
		err = rows.StructScan(webComponent)
		if err != nil {
			return nil, err
		}
		webComponentList = append(webComponentList, webComponent)
	}
	return webComponentList, nil
}

func (s *WebComponentDAO) DeleteComponent(id int) error {
	return nil
}
func (s *WebComponentDAO) UpdateComponent(id int) error {
	return nil
}
