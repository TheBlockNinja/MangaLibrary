package dao

import (
	"MangaLibrary/src/internal/dto"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BooksDAO struct {
	DB *sqlx.DB
}

func (s *BooksDAO) NewBook(books *dto.Books) error {
	if books.Name == "" {
		return fmt.Errorf("missing name")
	}
	if books.FilePath == "" {
		return fmt.Errorf("missing filepath")
	}
	insertStmt := "insert into manga_library.books(user_id,job_id,chapter,volume,name,description,meta_data,file_path,cover_img,pages) VALUES(?,?,?,?,?,?,?,?,?,?)"
	results, err := s.DB.Exec(insertStmt, books.UserID, books.JobId, books.Chapter, books.Volume, books.Name, books.Description, books.Metadata, books.FilePath, books.CoverImage, books.Pages)
	if err != nil {
		return err
	}
	id, _ := results.LastInsertId()
	books.ID = int(id)
	return nil
}

func (s *BooksDAO) GetBooksForName(name string, userId int) ([]*dto.Books, error) {
	var bookList []*dto.Books
	rows, err := s.DB.Queryx("select * from manga_library.books where name = ? and user_id = ? order by volume,chapter", name, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		book := &dto.Books{}
		err = rows.StructScan(book)
		if err != nil {
			return nil, err
		}
		bookList = append(bookList, book)
	}
	return bookList, nil
}

func (s *BooksDAO) GetAllUserBooks(userId int) ([]*dto.Books, error) {
	var bookList []*dto.Books
	rows, err := s.DB.Queryx("select * from manga_library.books where user_id = ? order by name,volume,chapter", userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		book := &dto.Books{}
		err = rows.StructScan(book)
		if err != nil {
			return nil, err
		}
		bookList = append(bookList, book)
	}
	return bookList, nil
}

func (s *BooksDAO) SearchForBooks() {

}

func (s *BooksDAO) GetDistinctBooks(userId int) ([]*dto.Books, error) {
	var bookList []*dto.Books
	rows, err := s.DB.Queryx("select * from manga_library.books where user_id = ? order by name,volume,chapter group by name", userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		book := &dto.Books{}
		err = rows.StructScan(book)
		if err != nil {
			return nil, err
		}
		bookList = append(bookList, book)
	}
	return bookList, nil
}
