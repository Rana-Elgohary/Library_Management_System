package routes

import (
	controllers "github.com/Pyramakerz/Library_Management_System/PKG/Controllers"
	"github.com/gofiber/fiber/v2"
)

// app *fiber.App ==> pointer to configure routes and middleware for your web application.
func Library_Management_System_Routes(app *fiber.App) {
	app.Get("/api/author", controllers.GetAllAuthors)
	app.Get("/api/author/:authorid", controllers.GetAuthorByID)
	app.Post("/api/author", controllers.CreateAuthor)
	app.Put("/api/author/:authorid", controllers.UpdateAuthor)
	app.Delete("/api/author/:authorid", controllers.DeleteAuthor)
	app.Delete("/api/author/softdelete/:authorid", controllers.SoftDeleteAuthor)

	app.Get("/api/book", controllers.GetAllBooks)
	app.Get("/api/book/:bookid", controllers.GetBookByID)
	app.Post("/api/book", controllers.CreateBook)
	app.Put("/api/book/:bookid", controllers.UpdateBook)
	app.Delete("/api/book/:bookid", controllers.DeleteBook)
	app.Delete("/api/book/softdelete/:bookid", controllers.SoftDeleteBook)
	app.Get("/api/book/search/:title", controllers.SearchBooksByTitle)
}
