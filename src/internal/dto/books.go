package dto

type Books struct {
	ID          int    `db:"id" json:"id"`
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
