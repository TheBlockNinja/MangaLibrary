package dto

type User struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	APIKey      string `db:"api_key"`
	CurrentJobs int    `db:"current_jobs"`
	MaxJobs     int    `db:"max_jobs"`
	Age         int    `db:"age"`
	LastLogin   string `db:"last_login"`
	IsActive    bool   `db:"is_active"`
	IsAdmin     bool   `db:"is_admin"`
}
