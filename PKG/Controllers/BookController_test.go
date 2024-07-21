package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	config "github.com/Pyramakerz/Library_Management_System/PKG/Config"
	models "github.com/Pyramakerz/Library_Management_System/PKG/Models"
	"github.com/stretchr/testify/assert"
)

func TestGetAllBooks(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	req := httptest.NewRequest(http.MethodGet, "/api/book", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetBookByID(t *testing.T) {
	app := SetupFiberApp()
	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	book := models.Book{
		Title:         "Sample Book",
		ISBN:          "1234567890",
		PublishedDate: time.Now(),
		AuthorID:      author.ID,
	}
	db.Create(&book)

	req := httptest.NewRequest(http.MethodGet, "/api/book/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateBook(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "Jane Doe", Email: "jane@example.com"}
	db.Create(&author)

	newBook := CreateBookRequest{
		Title:         "New Book",
		ISBN:          "0987654321",
		PublishedDate: time.Now(),
		AuthorID:      author.ID,
	}
	body, _ := json.Marshal(newBook)

	req := httptest.NewRequest(http.MethodPost, "/api/book", bytes.NewReader(body))
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUpdateBook(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	book := models.Book{
		Title:         "Old Book",
		ISBN:          "1234567890",
		PublishedDate: time.Now(),
		AuthorID:      author.ID,
	}
	db.Create(&book)

	updatedBook := CreateBookRequest{
		Title:         "Updated Book",
		ISBN:          "0987654321",
		PublishedDate: time.Now(),
		AuthorID:      author.ID,
	}
	body, _ := json.Marshal(updatedBook)

	req := httptest.NewRequest(http.MethodPut, "/api/book/1", bytes.NewReader(body))
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteBook(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	book := models.Book{
		Title:         "Book to Delete",
		ISBN:          "1234567890",
		PublishedDate: time.Now(),
		AuthorID:      author.ID,
	}
	db.Create(&book)

	req := httptest.NewRequest(http.MethodDelete, "/api/book/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSoftDeleteBook(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	book := models.Book{
		Title:         "Book to Soft Delete",
		ISBN:          "1234567890",
		PublishedDate: time.Now(),
		AuthorID:      author.ID,
	}
	db.Create(&book)

	req := httptest.NewRequest(http.MethodDelete, "/api/book/softdelete/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSearchBooksByTitle(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	book := models.Book{
		Title:         "Searchable Book",
		ISBN:          "1234567890",
		PublishedDate: time.Now(),
		AuthorID:      author.ID,
	}
	db.Create(&book)

	req := httptest.NewRequest(http.MethodGet, "/api/book/search/Searchable", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
