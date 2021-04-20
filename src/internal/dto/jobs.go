package dto

type Job struct {
	ID              int    `db:"id" json:"id"`
	UserID          int    `db:"user_id" json:"user_id"`
	SiteID          int    `db:"site_id" json:"site_id"`
	Name            string `db:"name" json:"name"`
	JobContext      string `db:"job_context" json:"job_context"`
	StartTime       string `db:"start_time" json:"start_time"`
	CurrentProgress int    `db:"current" json:"current"`
	TotalProgress   int    `db:"total" json:"total"`
	EstFinish       string `db:"est_finish" json:"est_finish"`
	Message         string `db:"message" json:"message"`
	CurrentJobData  string `db:"job_data" json:"job_data"`
}
