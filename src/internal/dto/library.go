package dto

type Library struct {
	UserID     int    `db:"user_id" json:"user_id"`
	BookID     int    `db:"book_id" json:"book_id"`
	Collection string `db:"collection" json:"collection"`
	Progress   int    `db:"progress" json:"progress"`
	Rating     int    `db:"rating" json:"rating"`
	Favorite   bool   `db:"favorite" json:"favorite"`
}
