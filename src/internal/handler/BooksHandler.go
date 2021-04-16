package handler

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"MangaLibrary/src/internal/pdf"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"go.uber.org/zap"
)

type BookRequest struct {
	UserID       int            `json:"user_id"`
	BookName     string         `json:"name"`
	BookID       int            `json:"book_id"`
	LibraryBooks []*dto.Library `json:"library_books"`
	Chapter      string         `json:"chapter"`
	Volume       string         `json:"volume"`
	Search       string         `json:"search"`
	Type         string         `json:"type"`
	Refresh      bool           `json:"refresh"`
	SendFile     bool           `json:"send_file"`
	Download     bool           `json:"download"`
	Image        string         `json:"image"`
}
type BookHandler struct {
	Site   *dao.SitesDAO
	Book   *dao.BooksDAO
	User   *dao.UsersDAO
	Logger *zap.Logger
}

func (h BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authentication-Key")
	u, err := h.User.GetUser(apiKey)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	bookReq := &BookRequest{}
	err = json.NewDecoder(r.Body).Decode(bookReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	books, err := h.Book.GetBooksForName(bookReq.BookName, u.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(books, "", w)
}

func (h BookHandler) GetBookDistinct(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authentication-Key")
	u, err := h.User.GetUser(apiKey)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	books, err := h.Book.GetDistinctBooks(u.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(books, "", w)
}

func (h BookHandler) GetBookSearch(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authentication-Key")
	u, err := h.User.GetUser(apiKey)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	bookReq := &BookRequest{}
	err = json.NewDecoder(r.Body).Decode(bookReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	fmt.Printf("searching for book %s\n user: %d\n", bookReq.Search, u.ID)
	books, err := h.Book.GetBooksSearchName(bookReq.Search, u.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(books, "", w)
}

func (h BookHandler) GetBookForID(w http.ResponseWriter, r *http.Request) {
	bookReq := &BookRequest{}
	err := json.NewDecoder(r.Body).Decode(bookReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	book, err := h.Book.GetBookForID(bookReq.BookID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(book, "", w)
}

func (h BookHandler) GetLibraryBooks(w http.ResponseWriter, r *http.Request) {
	bookReq := &BookRequest{}
	err := json.NewDecoder(r.Body).Decode(bookReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	books := []*dto.Books{}
	for _, l := range bookReq.LibraryBooks {
		book, err := h.Book.GetBookForID(l.BookID)
		if err != nil {
			SendError(err.Error(), w)
			return
		}
		books = append(books, book)
	}
	SendData(books, "", w)
}

func (h BookHandler) GetBookPDF(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if _, found := query["api_key"]; !found {
		SendError("did not find api_key", w)
		return
	}
	apiKey := query["api_key"][0]
	u, err := h.User.GetUser(apiKey)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	if _, found := query["book"]; !found {
		SendError("did not find book", w)
		return
	}
	books, err := h.Book.GetBooksForName(query["book"][0], u.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	if volume, found := query["volume"]; found {
		newBooks := []*dto.Books{}
		if len(volume) > 0 {
			for _, b := range books {
				if b.Volume == volume[0] {
					newBooks = append(newBooks, b)
				}
			}
			books = newBooks
		}
	}

	if chapter, found := query["chapter"]; found {
		if len(chapter) > 0 {
			for _, b := range books {
				if b.Chapter == chapter[0] {
					books = []*dto.Books{b}
					break
				}
			}
		}
	}
	formatType := "pdf"
	refresh := false
	download := false
	sendFile := true
	if fType, found := query["type"]; found {
		formatType = fType[0]
	}
	if ref, found := query["refresh"]; found {
		refresh, err = strconv.ParseBool(ref[0])
	}
	if down, found := query["download"]; found {
		download, err = strconv.ParseBool(down[0])
	}
	if down, found := query["send_file"]; found {
		sendFile, err = strconv.ParseBool(down[0])
	}
	output := []string{}
	switch formatType {
	case "pdf":
		for _, b := range books {
			baseFile := b.FilePath
			filename := fmt.Sprintf("%spdf/%s", baseFile, b.Name)
			if b.Chapter != "" {
				filename += " " + b.Chapter
			}
			if b.Volume != "" {
				filename += " " + b.Volume
			}
			err = os.MkdirAll(filename, os.ModePerm)
			if err != nil {
				SendError(err.Error(), w)
				return
			}
			filename += ".pdf"
			err = pdf.CreatePDFV2(b.FilePath, filename, refresh)
			if err != nil {
				h.Logger.Error("failed creating pdf", zap.Error(err))
				SendError(err.Error(), w)
				return
			}
			output = append(output, filename)
		}
	default:
		SendError("type not supported", w)
	}
	if sendFile {
		for _, f := range output {
			SendFile(f, download, w, r)
		}
	} else {
		SendData(output, "", w)
	}

}

func (h BookHandler) GetBookCover(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	image := ""
	if len(query["image"][0]) == 0 {
		apiKey := query["api_key"][0]
		u, err := h.User.GetUser(apiKey)
		if err != nil {
			SendError(err.Error(), w)
			return
		}
		bookName := query["book"][0]
		bookName, err = url.PathUnescape(bookName)
		if err != nil {
			SendError(err.Error(), w)
			return
		}
		books, err := h.Book.GetBooksForName(bookName, u.ID)
		if err != nil {
			SendError(err.Error(), w)
			return
		}

		if chapter, found := query["chapter"]; found {
			for _, b := range books {
				if b.Chapter == chapter[0] {
					SendFile(b.CoverImage, false, w, r)
					return
				}
			}
		}
		image = books[0].CoverImage
	} else {
		image = query["image"][0]
	}
	SendFile(image, false, w, r)
	return

}

func (h BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authentication-Key")
	u, err := h.User.GetUser(apiKey)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	books, err := h.Book.GetAllUserBooks(u.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(books, "", w)
}
