package dto

type WebComponent struct {
	ID             int    `db:"id" json:"id"`
	SiteID         int    `db:"site_id"  json:"site_id"`
	Name           string `db:"name"  json:"name"`
	Tag            string `db:"tag"  json:"tag"`
	Attribute      string `db:"attribute"  json:"attribute"`
	Value          string `db:"value"  json:"value"`
	IsLink         bool   `db:"is_link"  json:"is_link"`
	IsDownload     bool   `db:"is_download"  json:"is_download"`
	LinkAttributes string `db:"link_attributes"  json:"link_attributes"`
	ElementData    string `db:"element_data"  json:"element_data"`
	Parent         int    `db:"parent"  json:"parent"`
	Delay          int    `db:"delay"  json:"delay"`
	MetaData       string `db:"meta_data"  json:"meta_data"`
	Reverse        bool   `db:"reverse" json:"reverse"`
}
