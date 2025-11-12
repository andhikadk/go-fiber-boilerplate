package handlers

import (
	"strconv"

	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/internal/utils"
	pkgUtils "go-fiber-boilerplate/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetBooks godoc
//
//	@Summary		Get all books
//	@Description	Retrieve all books with pagination
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int	false	"Page number (default: 1)"
//	@Param			limit	query		int	false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	models.PaginatedResponse{data=[]models.Book}	"Books retrieved successfully"
//	@Failure		401		{object}	models.APIResponse							"Unauthorized"
//	@Failure		500		{object}	models.APIResponse							"Failed to fetch books"
//	@Router			/api/books [get]
func GetBooks(c *fiber.Ctx) error {
	// Get pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get books from service
	bookService := services.NewBookService(database.GetDB())
	books, total, err := bookService.GetAllBooks(page, limit)
	if err != nil {
		utils.ErrorLogger.Printf("[GetBooks] Failed to fetch books: %v", err)
		return pkgUtils.InternalErrorResponse(c, "Failed to fetch books")
	}

	utils.InfoLogger.Printf("[GetBooks] Retrieved %d books (page: %d, limit: %d)", len(books), page, limit)
	return pkgUtils.PaginatedResponse(c, "Books retrieved successfully", books, page, limit, total)
}

// GetBook godoc
//
//	@Summary		Get a book by ID
//	@Description	Retrieve a specific book by its ID
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Book ID"
//	@Success		200	{object}	models.APIResponse{data=models.Book}	"Book retrieved successfully"
//	@Failure		400	{object}	models.APIResponse						"Invalid book ID"
//	@Failure		401	{object}	models.APIResponse						"Unauthorized"
//	@Failure		404	{object}	models.APIResponse						"Book not found"
//	@Router			/api/books/{id} [get]
func GetBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		utils.ErrorLogger.Printf("[GetBook] Invalid book ID: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid book ID")
	}

	bookService := services.NewBookService(database.GetDB())
	book, err := bookService.GetBookByID(uint(id))
	if err != nil {
		utils.ErrorLogger.Printf("[GetBook] Book not found (ID: %d): %v", id, err)
		return pkgUtils.NotFoundResponse(c, err.Error())
	}

	utils.InfoLogger.Printf("[GetBook] Book retrieved successfully (ID: %d)", id)
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Book retrieved successfully", book)
}

// CreateBook godoc
//
//	@Summary		Create a new book
//	@Description	Create a new book with title, author, year, and ISBN
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateBookRequest		true	"Book creation data"
//	@Success		201		{object}	models.APIResponse{data=models.Book}	"Book created successfully"
//	@Failure		400		{object}	models.APIResponse			"Invalid request or validation error"
//	@Failure		401		{object}	models.APIResponse			"Unauthorized"
//	@Failure		500		{object}	models.APIResponse			"Failed to create book"
//	@Router			/api/books [post]
func CreateBook(c *fiber.Ctx) error {
	var req dto.CreateBookRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[CreateBook] Failed to parse request body: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request using DTO's self-validation
	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[CreateBook] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	// Create book
	bookService := services.NewBookService(database.GetDB())
	book, err := bookService.CreateBook(&req)
	if err != nil {
		utils.ErrorLogger.Printf("[CreateBook] Failed to create book: %v", err)
		return pkgUtils.InternalErrorResponse(c, "Failed to create book")
	}

	utils.InfoLogger.Printf("[CreateBook] Book created successfully (ID: %d, Title: %s)", book.ID, book.Title)
	return pkgUtils.CreatedResponse(c, "Book created successfully", book)
}

// UpdateBook godoc
//
//	@Summary		Update a book
//	@Description	Update an existing book's information
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Book ID"
//	@Param			request	body		dto.UpdateBookRequest	true	"Book update data"
//	@Success		200		{object}	models.APIResponse{data=models.Book}	"Book updated successfully"
//	@Failure		400		{object}	models.APIResponse		"Invalid request or validation error"
//	@Failure		401		{object}	models.APIResponse		"Unauthorized"
//	@Failure		404		{object}	models.APIResponse		"Book not found"
//	@Router			/api/books/{id} [put]
func UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		utils.ErrorLogger.Printf("[UpdateBook] Invalid book ID: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid book ID")
	}

	var req dto.UpdateBookRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[UpdateBook] Failed to parse request body: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request using DTO's self-validation
	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[UpdateBook] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	// Update book
	bookService := services.NewBookService(database.GetDB())
	book, err := bookService.UpdateBook(uint(id), &req)
	if err != nil {
		utils.ErrorLogger.Printf("[UpdateBook] Failed to update book (ID: %d): %v", id, err)
		return pkgUtils.NotFoundResponse(c, "Book not found")
	}

	utils.InfoLogger.Printf("[UpdateBook] Book updated successfully (ID: %d)", id)
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Book updated successfully", book)
}

// DeleteBook godoc
//
//	@Summary		Delete a book
//	@Description	Delete a book by its ID (soft delete)
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Book ID"
//	@Success		200	{object}	models.APIResponse	"Book deleted successfully"
//	@Failure		400	{object}	models.APIResponse	"Invalid book ID"
//	@Failure		401	{object}	models.APIResponse	"Unauthorized"
//	@Failure		404	{object}	models.APIResponse	"Book not found"
//	@Router			/api/books/{id} [delete]
func DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		utils.ErrorLogger.Printf("[DeleteBook] Invalid book ID: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid book ID")
	}

	// Delete book
	bookService := services.NewBookService(database.GetDB())
	if err := bookService.DeleteBook(uint(id)); err != nil {
		utils.ErrorLogger.Printf("[DeleteBook] Failed to delete book (ID: %d): %v", id, err)
		return pkgUtils.NotFoundResponse(c, "Book not found")
	}

	utils.InfoLogger.Printf("[DeleteBook] Book deleted successfully (ID: %d)", id)
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Book deleted successfully", nil)
}

// SearchBooks godoc
//
//	@Summary		Search books
//	@Description	Search for books by title or author
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			q	query		string	true	"Search query"
//	@Success		200	{object}	models.APIResponse{data=[]models.Book}	"Search results"
//	@Failure		400	{object}	models.APIResponse						"Search query is required"
//	@Failure		401	{object}	models.APIResponse						"Unauthorized"
//	@Failure		500	{object}	models.APIResponse						"Search failed"
//	@Router			/api/books/search [get]
func SearchBooks(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		utils.ErrorLogger.Printf("[SearchBooks] Empty search query")
		return pkgUtils.BadRequestResponse(c, "Search query is required")
	}

	bookService := services.NewBookService(database.GetDB())
	books, err := bookService.SearchBooks(query)
	if err != nil {
		utils.ErrorLogger.Printf("[SearchBooks] Search failed for query '%s': %v", query, err)
		return pkgUtils.InternalErrorResponse(c, "Search failed")
	}

	utils.InfoLogger.Printf("[SearchBooks] Search completed for query '%s', found %d books", query, len(books))
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Search results", books)
}
