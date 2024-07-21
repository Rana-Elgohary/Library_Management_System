package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/Pyramakerz/Library_Management_System/PKG/Config"
	models "github.com/Pyramakerz/Library_Management_System/PKG/Models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func SetupFiberApp() *fiber.App {
	app := fiber.New()
	config.Connect()
	db := config.GetDB()

	db.AutoMigrate(&models.Author{})

	app.Get("/api/author", GetAllAuthors)
	app.Get("/api/author/:authorid", GetAuthorByID)
	app.Post("/api/author", CreateAuthor)
	app.Put("/api/author/:authorid", UpdateAuthor)
	app.Delete("/api/author/:authorid", DeleteAuthor)
	app.Delete("/api/author/softdelete/:authorid", SoftDeleteAuthor)

	db.AutoMigrate(&models.Author{}, &models.Book{})

	app.Get("/api/book", GetAllBooks)
	app.Get("/api/book/:bookid", GetBookByID)
	app.Post("/api/book", CreateBook)
	app.Put("/api/book/:bookid", UpdateBook)
	app.Delete("/api/book/:bookid", DeleteBook)
	app.Delete("/api/book/softdelete/:bookid", SoftDeleteBook)
	app.Get("/api/book/search/:title", SearchBooksByTitle)
	return app
}

// But it will delete all what is inside the table
func CleanDB(db *gorm.DB) {
	db.Exec("DELETE FROM books")
	db.Exec("DELETE FROM authors")
	db.Exec("ALTER TABLE books AUTO_INCREMENT = 1")
	db.Exec("ALTER TABLE authors AUTO_INCREMENT = 1")
}

func TestGetAllAuthors(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	req := httptest.NewRequest(http.MethodGet, "/api/author", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAuthorByID(t *testing.T) {
	app := SetupFiberApp()
	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	req := httptest.NewRequest(http.MethodGet, "/api/author/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateAuthor(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "Jane Doe", Email: "jane@example.com"}
	body, _ := json.Marshal(author)

	req := httptest.NewRequest(http.MethodPost, "/api/author", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUpdateAuthor(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	updatedAuthor := models.Author{Name: "John Updated", Email: "johnupdated@example.com"}
	body, _ := json.Marshal(updatedAuthor)

	req := httptest.NewRequest(http.MethodPut, "/api/author/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteAuthor(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	req := httptest.NewRequest(http.MethodDelete, "/api/author/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSoftDeleteAuthor(t *testing.T) {
	app := SetupFiberApp()

	db := config.GetDB()
	defer CleanDB(db)

	author := models.Author{Name: "John Doe", Email: "john@example.com"}
	db.Create(&author)

	req := httptest.NewRequest(http.MethodDelete, "/api/author/softdelete/1", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
