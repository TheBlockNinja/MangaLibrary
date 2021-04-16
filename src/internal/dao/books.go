package dao

import (
	"MangaLibrary/src/internal/dto"
	"database/sql"
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
	existingBooks, err := s.GetBooksForName(books.Name, books.UserID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, b := range existingBooks {
		if books.Volume == b.Volume && books.Chapter == b.Chapter && books.FilePath == b.FilePath {
			fmt.Printf("%s ID %d: book already exists\n", books.Name, b.ID)
			return nil
		}
	}
	if books.FilePath == "" {
		return fmt.Errorf("missing filepath")
	}
	insertStmt := "insert into manga_library.books(user_id,job_id,chapter,volume,name,description,meta_data,file_path,cover_img,pages,site_id,is_public,views,downloads) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	results, err := s.DB.Exec(insertStmt, books.UserID, books.JobId, books.Chapter, books.Volume, books.Name, books.Description, books.Metadata, books.FilePath, books.CoverImage, books.Pages, books.SiteID, books.Public, books.Views, books.Downloads)
	if err != nil {
		return err
	}
	id, _ := results.LastInsertId()
	books.ID = int(id)
	return nil
}

func (s *BooksDAO) GetBooksForName(name string, userId int) ([]*dto.Books, error) {
	var bookList []*dto.Books
	rows, err := s.DB.Queryx("select * from manga_library.books where `name` = ? and user_id = ? order by id desc", name, userId)
	if err == sql.ErrNoRows {
		return bookList, nil
	}
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
func (s *BooksDAO) GetBookForID(id int) (*dto.Books, error) {
	book := &dto.Books{}
	err := s.DB.Get(book, "select * from manga_library.books where id = ?", id)
	if err != nil {
		return book, err
	}
	return book, nil
}

func (s *BooksDAO) UpdateBook(book *dto.Books, userId int) error {
	query := "update manga_library.books set cover_img = ?, chapter=?, volume=?,name=?,description=?,is_public=?,views=?,downloads=? where id = ? and user = ?"
	_, err := s.DB.Exec(query, book.CoverImage, book.Chapter, book.Volume, book.Name, book.Description, book.Public, book.Views, book.Downloads, book.ID, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *BooksDAO) GetBooksSearchName(name string, userId int) ([]*dto.Books, error) {
	var bookList []*dto.Books
	query := "select * from manga_library.books where `name` like ? "
	rows, err := s.DB.Queryx(query, "%"+name+"%s")
	if err == sql.ErrNoRows {
		return bookList, nil
	}
	if err != nil {
		return nil, err
	}
	fmt.Printf("err %s\n", rows.Err())
	for rows.Next() {
		fmt.Printf("getting row...")
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
	rows, err := s.DB.Queryx("select * from manga_library.books where user_id = ? or is_public=1 order by id desc", userId)
	if err == sql.ErrNoRows {
		return bookList, nil
	}
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
	rows, err := s.DB.Queryx("select * from manga_library.books where user_id = ? or is_public = 1  order by id", userId)
	if err == sql.ErrNoRows {
		return bookList, nil
	}
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

	output := []*dto.Books{}
	tmpMap := map[string]string{}
	for _, b := range bookList {
		if _, found := tmpMap[b.Name]; !found {
			output = append(output, b)
			tmpMap[b.Name] = ""
		}
	}
	return output, nil
}
