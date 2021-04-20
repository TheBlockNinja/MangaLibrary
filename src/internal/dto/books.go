package dto

type Books struct {
	ID          int    `db:"id" json:"id"`
	SiteID      int    `db:"site_id" json:"site_id"`
	SiteName    string `db:"site_name" json:"site_name"`
	Public      bool   `db:"is_public" json:"is_public"`
	Views       int    `db:"views" json:"views"`
	Downloads   int    `db:"downloads" json:"downloads"`
	UserID      int    `db:"user_id" json:"user_id"`
	JobId       int    `db:"job_id" json:"job_id"`
	Chapter     string `db:"chapter" json:"chapter"`
	Volume      string `db:"volume" json:"volume"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Metadata    string `db:"meta_data" json:"meta_data"`
	FilePath    string `db:"file_path" json:"file_path"`
	CoverImage  string `db:"cover_img" json:"cover_img"`
	Pages       int    `db:"pages" json:"pages"`
}
