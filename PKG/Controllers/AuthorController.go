package controllers

import (
	"errors"
	"fmt"

	"github.com/MakMoinee/go-mith/pkg/email"
	config "github.com/Pyramakerz/Library_Management_System/PKG/Config"
	models "github.com/Pyramakerz/Library_Management_System/PKG/Models"
	utils "github.com/Pyramakerz/Library_Management_System/PKG/Utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ----------------------------------------------------------------------------------------------------------------------------------

var db = config.GetDB()

func ensureDB() {
	if db == nil {
		db = config.GetDB()
		if db == nil {
			panic("Database connection is not initialized")
		}
	}
}

// ----------------------------------------------------------------------------------------------------------------------------------

// GetAllAuthors godoc
// @Summary      Get all authors
// @Description  Get a list of all authors
// @Tags         authors
// @Accept       json
// @Produce      json
// @Success      200  {object}  any
// @Failure      500  {object}  any
// @Router       /api/author [get]
func GetAllAuthors(c *fiber.Ctx) error {
	ensureDB()

	var authors []models.Author
	if err := db.Find(&authors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch authors",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  authors,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// GetAuthorByID godoc
// @Summary      Get author by ID
// @Description  Get a specific author by their ID
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        authorid  path  string  true  "Author ID"
// @Success      200  {object}  models.Author
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      500  {object}  any
// @Router       /api/author/{authorid} [get]
func GetAuthorByID(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("authorid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Author ID",
		})
	}

	var author models.Author
	if err := db.First(&author, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get author",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  author,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// CreateAuthor godoc
// @Summary      Create a new author
// @Description  Create a new author with the provided information
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        author  body  models.Author  true  "Author data"
// @Success      201  {object}  models.Author
// @Failure      400  {object}  any
// @Failure      409  {object}  any
// @Failure      500  {object}  any
// @Router       /api/author [post]
func CreateAuthor(c *fiber.Ctx) error {
	ensureDB()

	author := models.Author{}

	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Cannot parse JSON",
		})
	}

	if author.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Name is required",
		})
	}

	if author.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Email is required",
		})
	}

	if !utils.IsValidEmail(author.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid email format",
		})
	}

	// Check if email already exists
	var existingAuthor models.Author
	if err := db.Where("email = ?", author.Email).First(&existingAuthor).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Email already exists",
		})
	}

	// Check if ID already exists
	if author.ID != 0 { // Assuming ID is an auto-increment field and should not be provided
		if err := db.First(&existingAuthor, author.ID).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "ID already exists",
			})
		}
	}

	// Create() ==> already make the save operation
	if err := db.Create(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create author",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"data":  author,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// UpdateAuthor godoc
// @Summary      Update an existing author
// @Description  Update an existing author's information
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        authorid  path  string  true  "Author ID"
// @Param        author    body  models.Author  true  "Updated author data"
// @Success      200  {object}  models.Author
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      409  {object}  any
// @Failure      500  {object}  any
// @Router       /api/author/{authorid} [put]
func UpdateAuthor(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("authorid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Author ID",
		})
	}

	var existingAuthor models.Author
	if err := db.First(&existingAuthor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to find author",
		})
	}

	var updatedAuthor models.Author

	if err := c.BodyParser(&updatedAuthor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Cannot parse JSON",
		})
	}

	if !utils.IsValidEmail(updatedAuthor.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid email format",
		})
	}

	// Check if the new email already exists for a different author
	// <> ==> not equal
	var conflictingAuthorEmail models.Author
	if err := db.Where("email = ? AND id <> ?", updatedAuthor.Email, id).First(&conflictingAuthorEmail).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Email already exists",
		})
	}

	if existingAuthor.Email != updatedAuthor.Email {
		//  SMTP (Simple Mail Transfer Protocol) 587 (its port number)
		emailService := email.NewEmailService(587, "smtp.gmail.com", "ranamelgohary2000@gmail.com", "inoxwuyyzogdvrfs")
		isSent, err := emailService.SendEmail("r.elgohary2000@gmail.com", "email notification", "Author information updated")
		if err != nil {
			fmt.Println(err)
		}
		if isSent {
			fmt.Println("Done")
		}
		fmt.Printf("Sending email notification to %s: Author information updated", updatedAuthor.Email)
	}

	existingAuthor.Name = updatedAuthor.Name
	existingAuthor.Email = updatedAuthor.Email

	if err := db.Save(&existingAuthor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update author",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  existingAuthor,
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// DeleteAuthor godoc
// @Summary      Delete an author
// @Description  Permanently delete an author by their ID
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        authorid  path  string  true  "Author ID"
// @Success      200  {object}  any
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      500  {object}  any
// @Router       /api/author/{authorid} [delete]
func DeleteAuthor(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("authorid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Author ID",
		})
	}

	var author models.Author

	// db.Unscoped() ==> for Normal delete
	if err := db.Unscoped().First(&author, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete author",
		})
	}

	// Delete() ==> already make the save operation, GORM will ignore the DeletedAt field and perform the operation as if the record is not soft-deleted
	if err := db.Unscoped().Delete(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete author",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Author deleted successfully",
	})
}

// ----------------------------------------------------------------------------------------------------------------------------------

// SoftDeleteAuthor godoc
// @Summary      Soft delete an author
// @Description  Soft delete an author by their ID (sets the deleted_at timestamp)
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        authorid  path  string  true  "Author ID"
// @Success      200  {object}  any
// @Failure      400  {object}  any
// @Failure      404  {object}  any
// @Failure      500  {object}  any
// @Router       /api/author/softdelete/{authorid} [delete]
func SoftDeleteAuthor(c *fiber.Ctx) error {
	ensureDB()

	id := c.Params("authorid")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Enter Author ID",
		})
	}

	var author models.Author
	if err := db.First(&author, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to find author",
		})
	}

	if err := db.Delete(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to soft delete author",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Author soft deleted successfully",
	})
}
