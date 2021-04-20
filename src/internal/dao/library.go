package dao

import (
	"MangaLibrary/src/internal/dto"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type LibraryDAO struct {
	DB *sqlx.DB
}

func (s *LibraryDAO) NewElement(library *dto.Library) error {
	insertStmt := "insert into manga_library.library(user_id,book_id,collection,progress,rating,favorite) VALUES(?,?,?,?,?,?)"
	_, err := s.DB.Exec(insertStmt, library.UserID, library.BookID, library.Collection, library.Progress, library.Rating, library.Favorite)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryDAO) UpdateItem(library *dto.Library) error {
	_, err := s.DB.Exec("update manga_library.library set collection = ?, progress = ?, rating = ?,favorite = ? where user_id = ? and book_id = ?",
		library.Collection, library.Progress, library.Rating, library.Favorite, library.UserID, library.BookID)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryDAO) DeleteItem(library *dto.Library) error {
	_, err := s.DB.Exec("Delete from manga_library.library where user_id = ? and book_id = ?", library.UserID, library.BookID)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryDAO) GetItemsForCollection(collection string, user int) ([]*dto.Library, error) {
	libraryList := []*dto.Library{}

	rows, err := s.DB.Queryx("select * from manga_library.library where user_id = ? and collection = ?", user, collection)
	if err == sql.ErrNoRows {
		return libraryList, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		webComponent := &dto.Library{}
		err = rows.StructScan(webComponent)
		if err != nil {
			return nil, err
		}
		libraryList = append(libraryList, webComponent)
	}
	return libraryList, nil
}

func (s *LibraryDAO) GetCollections(user int) ([]string, error) {
	collectionList := []string{}
	rows, err := s.DB.Queryx("select distinct(collection) from manga_library.library where user_id = ?", user)
	if err == sql.ErrNoRows {
		return collectionList, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		webComponent := ""
		err = rows.Scan(&webComponent)
		if err != nil {
			return nil, err
		}
		collectionList = append(collectionList, webComponent)
	}
	return collectionList, nil
}

func (s *LibraryDAO) GetLibraryItem(user, book int) (*dto.Library, error) {
	library := &dto.Library{}
	err := s.DB.Get(library, "select * from manga_library.library where user_id = ? and book_id = ?", user, book)
	if err != nil {
		return nil, err
	}
	return library, nil
}
