package service

import (
	"bookservice/cache"
	"bookservice/models"

	"github.com/google/uuid"
)

type BookService struct {
	cache cache.IRedisClient
}

type IBookService interface {
	CreateBook(data *models.CreateBookRequestModel) error
	DeleteBookById(bookID uuid.UUID) error
	GetBookById(bookID uuid.UUID) (*models.BookModel, error)
	GetAllBooks(queries map[string]string) ([]models.BookModel, error)
	UpdatePriceById(bookID uuid.UUID, newPrice string) error
}

func NewBookService(cache cache.IRedisClient) IBookService {
	return &BookService{
		cache: cache,
	}
}

func (bookService *BookService) CreateBook(data *models.CreateBookRequestModel) error {
	return bookService.cache.CreateBooks(data)
}

func (bookService *BookService) DeleteBookById(bookID uuid.UUID) error {
	return bookService.cache.DeleteBookById(bookID)
}

func (bookService *BookService) GetBookById(bookID uuid.UUID) (*models.BookModel, error) {
	return bookService.cache.GetBookById(bookID)
}

func (bookService *BookService) GetAllBooks(queries map[string]string) ([]models.BookModel, error) {
	return bookService.cache.GetAllBooks(queries)
}

func (bookService *BookService) UpdatePriceById(bookID uuid.UUID, newPrice string) error {
	return bookService.cache.UpdatePriceById(bookID, newPrice)
}
