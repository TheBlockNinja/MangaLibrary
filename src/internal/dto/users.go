package dto

type User struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Email       string `db:"email" json:"email"`
	Password    string `db:"password" json:"password"`
	APIKey      string `db:"api_key" json:"api_key"`
	CurrentJobs int    `db:"current_jobs" json:"current_jobs"`
	MaxJobs     int    `db:"max_jobs" json:"max_jobs"`
	Age         int    `db:"age" json:"age"`
	LastLogin   string `db:"last_login" json:"last_login"`
	IsActive    bool   `db:"is_active" json:"is_active"`
	IsAdmin     bool   `db:"is_admin" json:"is_admin"`
}
