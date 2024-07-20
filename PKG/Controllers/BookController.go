package controllers

import (
	"errors"
	"time"

	models "github.com/Pyramakerz/Library_Management_System/PKG/Models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateBookRequest struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	ISBN          string    `json:"isbn"`
	PublishedDate time.Time `json:"publishedDate"`
	AuthorID      uint      `json:"authorID"`
}

// ----------------------------------------------------------------------------------------------------------------------------------

// GetAllBooks godoc
// @Summary      Get all books
// @Description  Get a list of all books, including their authors
// @Tags         books
// @Accept       json
// @Produce      json
// @Success      200  {object}  any
// @Failure      500  {object}  any
// @Router       /api/book [get]
func GetAllBooks(c *fiber.Ctx) error {
	ensureDB()

	var Books []models.Book
	if err := db.Preload("Author").Find(&Books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch books",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  Books,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// GetBookByID godoc
// @Summary      Get book by ID
// @Description  Get a specific book by its ID, including its author details
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        bookid  path  string  true  "Book ID"
// @Success      200  {object}  models.Book
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      500  {object}  any
// @Router       /api/book/{bookid} [get]
func GetBookByID(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("bookid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Book ID",
		})
	}

	var book models.Book
	if err := db.Preload("Author").First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get Book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  book,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// CreateBook godoc
// @Summary      Create a new book
// @Description  Create a new book with the provided information, including author details
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body  CreateBookRequest  true  "Book data"
// @Success      201  {object}  models.Book
// @Failure      400  {object}  any
// @Failure      409  {object}  any
// @Failure      500  {object}  any
// @Router       /api/book [post]
func CreateBook(c *fiber.Ctx) error {
	ensureDB()

	book := CreateBookRequest{}

	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Cannot parse JSON",
		})
	}

	if book.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Title is required",
		})
	}

	if book.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ISBN is required",
		})
	}

	if book.PublishedDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Published date is required",
		})
	}

	// Check if ISBN already exists
	var existingBook models.Book
	if err := db.Where("isbn = ?", book.ISBN).First(&existingBook).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "ISBN already exists",
		})
	}

	// Check if the author exists
	var author models.Author

	if err := db.First(&author, book.AuthorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to find author",
		})
	}

	// Only check if ID is provided
	// For not soft deleted only
	if book.ID != 0 {
		if err := db.First(&existingBook, book.ID).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "ID already exists",
			})
		}
	}

	NewBook := models.Book{
		Title:         book.Title,
		ISBN:          book.ISBN,
		PublishedDate: book.PublishedDate,
		AuthorID:      book.AuthorID,
		Author:        author,
	}

	if err := db.Create(&NewBook).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create book",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"data":  NewBook,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// UpdateBook godoc
// @Summary      Update an existing book
// @Description  Update an existing book's information, including author details
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        bookid  path  string  true  "Book ID"
// @Param        book    body  CreateBookRequest  true  "Updated book data"
// @Success      200  {object}  models.Book
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      409  {object}  any
// @Failure      500  {object}  any
// @Router       /api/book/{bookid} [put]
func UpdateBook(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("bookid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Book ID",
		})
	}

	var existingBook models.Book
	if err := db.Preload("Author").First(&existingBook, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to find book",
		})
	}

	var updatedBook CreateBookRequest
	if err := c.BodyParser(&updatedBook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Cannot parse JSON",
		})
	}

	if updatedBook.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Title is required",
		})
	}

	if updatedBook.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ISBN is required",
		})
	}

	if updatedBook.PublishedDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Published date is required",
		})
	}

	// Check if ISBN already exists and is not the current book's ISBN
	if updatedBook.ISBN != existingBook.ISBN {
		var existingBookWithNewISBN models.Book
		if err := db.Where("isbn = ?", updatedBook.ISBN).First(&existingBookWithNewISBN).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "ISBN already exists",
			})
		}
	}

	// Check if the author exists
	var author models.Author
	if err := db.First(&author, updatedBook.AuthorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to find author",
		})
	}

	// Update the book
	existingBook.Title = updatedBook.Title
	existingBook.ISBN = updatedBook.ISBN
	existingBook.PublishedDate = updatedBook.PublishedDate
	existingBook.AuthorID = updatedBook.AuthorID
	existingBook.Author = author

	if err := db.Save(&existingBook).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  existingBook,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// DeleteBook godoc
// @Summary      Delete a book
// @Description  Permanently delete a book by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        bookid  path  string  true  "Book ID"
// @Success      200  {object}  any
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      500  {object}  any
// @Router       /api/book/{bookid} [delete]
func DeleteBook(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("bookid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Book ID",
		})
	}

	var book models.Book

	// db.Unscoped() ==> for Normal delete
	if err := db.Unscoped().First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete book",
		})
	}

	// Delete() ==> already make the save operation, GORM will ignore the DeletedAt field and perform the operation as if the record is not soft-deleted
	if err := db.Unscoped().Delete(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Book deleted successfully",
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// SoftDeleteBook godoc
// @Summary      Soft delete a book
// @Description  Soft delete a book by its ID (sets the deleted_at timestamp)
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        bookid  path  string  true  "Book ID"
// @Success      200  {object}  any
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      500  {object}  any
// @Router       /api/book/softdelete/{bookid} [delete]
func SoftDeleteBook(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("bookid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Book ID",
		})
	}

	var book models.Book
	if err := db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to find book",
		})
	}

	if err := db.Delete(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to soft delete book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Book soft deleted successfully",
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// SearchBooksByTitle godoc
// @Summary      Search books by title
// @Description  Search for books based on a partial or full title match
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        title  path  string  true  "Book title"
// @Success      200  {object}  any
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      500  {object}  any
// @Router       /api/book/search/{title} [get]
func SearchBooksByTitle(c *fiber.Ctx) error {
	ensureDB()

	title := c.Params("title")

	if title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Title is required",
		})
	}

	var books []models.Book
	if err := db.Preload("Author").Where("title LIKE ?", "%"+title+"%").Find(&books).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "No books found with matching title",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to search books",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  books,
	})
}
