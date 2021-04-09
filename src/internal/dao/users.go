package dao

import (
	"MangaLibrary/src/internal/dto"
	"math/rand"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type UsersDAO struct {
	DB *sqlx.DB
}

func (s *UsersDAO) NewUser(user *dto.User) error {
	user.APIKey = GenerateAPIKey(20)
	insertStmt := "insert into manga_library.users(name,email,password,api_key,max_jobs,age,is_active,is_admin) VALUES(?,?,?,?,?,?,?,?)"
	results, err := s.DB.Exec(insertStmt, user.Name, user.Email, user.Password, user.APIKey, user.MaxJobs, user.Age, user.IsActive, user.IsAdmin)
	if err != nil {
		return err
	}
	id, _ := results.LastInsertId()
	user.ID = int(id)
	return nil
}
func (s *UsersDAO) UpdateUserAPIKey(user *dto.User) error {
	newAPIKey := GenerateAPIKey(20)
	_, err := s.DB.Exec("update manga_library.users set api_key = ? where id = ? and api_key = ?", newAPIKey, user.ID, user.APIKey)
	if err != nil {
		return err
	}
	user.APIKey = newAPIKey
	return nil
}
func GenerateAPIKey(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func (s *UsersDAO) GetUser(apiKey string) (*dto.User, error) {
	user := &dto.User{}
	err := s.DB.Get(user, "select * from manga_library.users where api_key = ?", apiKey)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UsersDAO) GetUserForName(name, password string) (*dto.User, error) {
	user := &dto.User{}
	err := s.DB.Get(user, "select * from manga_library.users where name = ? and password = ?", name, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
