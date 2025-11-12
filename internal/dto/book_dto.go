package dto

import (
	"errors"
	"strings"
	"time"
)

// CreateBookRequest is the request body for creating a new book
type CreateBookRequest struct {
	Title  string `json:"title" example:"The Go Programming Language"`
	Author string `json:"author" example:"Alan A. A. Donovan"`
	Year   int    `json:"year" example:"2015"`
	ISBN   string `json:"isbn" example:"978-0134190440"`
}

// Validate validates the CreateBookRequest
func (r *CreateBookRequest) Validate() error {
	// Validate Title
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required and cannot be empty")
	}
	if len(r.Title) < 2 {
		return errors.New("title must be at least 2 characters")
	}
	if len(r.Title) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}

	// Validate Author
	if strings.TrimSpace(r.Author) == "" {
		return errors.New("author is required and cannot be empty")
	}
	if len(r.Author) < 2 {
		return errors.New("author must be at least 2 characters")
	}
	if len(r.Author) > 255 {
		return errors.New("author cannot exceed 255 characters")
	}

	// Validate Year
	if r.Year < 1000 || r.Year > 9999 {
		return errors.New("year must be a valid 4-digit year (1000-9999)")
	}

	// Validate ISBN
	if strings.TrimSpace(r.ISBN) == "" {
		return errors.New("isbn is required and cannot be empty")
	}
	if len(r.ISBN) > 20 {
		return errors.New("isbn cannot exceed 20 characters")
	}

	return nil
}

// UpdateBookRequest is the request body for updating an existing book
type UpdateBookRequest struct {
	Title  *string `json:"title,omitempty" example:"The Go Programming Language (Updated)"`
	Author *string `json:"author,omitempty" example:"Alan A. A. Donovan & Brian W. Kernighan"`
	Year   *int    `json:"year,omitempty" example:"2016"`
	ISBN   *string `json:"isbn,omitempty" example:"978-0134190440"`
}

// Validate validates the UpdateBookRequest
func (r *UpdateBookRequest) Validate() error {
	// Validate Title if provided
	if r.Title != nil {
		if strings.TrimSpace(*r.Title) == "" {
			return errors.New("title cannot be empty")
		}
		if len(*r.Title) < 2 {
			return errors.New("title must be at least 2 characters")
		}
		if len(*r.Title) > 255 {
			return errors.New("title cannot exceed 255 characters")
		}
	}

	// Validate Author if provided
	if r.Author != nil {
		if strings.TrimSpace(*r.Author) == "" {
			return errors.New("author cannot be empty")
		}
		if len(*r.Author) < 2 {
			return errors.New("author must be at least 2 characters")
		}
		if len(*r.Author) > 255 {
			return errors.New("author cannot exceed 255 characters")
		}
	}

	// Validate Year if provided
	if r.Year != nil {
		if *r.Year < 1000 || *r.Year > 9999 {
			return errors.New("year must be a valid 4-digit year (1000-9999)")
		}
	}

	// Validate ISBN if provided
	if r.ISBN != nil {
		if strings.TrimSpace(*r.ISBN) == "" {
			return errors.New("isbn cannot be empty")
		}
		if len(*r.ISBN) > 20 {
			return errors.New("isbn cannot exceed 20 characters")
		}
	}

	return nil
}

// BookResponse is the response for book data
type BookResponse struct {
	ID        uint      `json:"id" example:"1"`
	Title     string    `json:"title" example:"The Go Programming Language"`
	Author    string    `json:"author" example:"Alan A. A. Donovan"`
	Year      int       `json:"year" example:"2015"`
	ISBN      string    `json:"isbn" example:"978-0134190440"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
