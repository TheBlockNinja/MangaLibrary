package dto

type Job struct {
	ID              int    `db:"id"`
	UserID          int    `db:"user_id"`
	SiteID          int    `db:"site_id"`
	Name            string `db:"name"`
	JobContext      string `db:"job_context"`
	StartTime       string `db:"start_time"`
	CurrentProgress int    `db:"current"`
	TotalProgress   int    `db:"total"`
	EstFinish       string `db:"est_finish"`
	Message         string `db:"message"`
	CurrentJobData  string `db:"job_data"`
}
