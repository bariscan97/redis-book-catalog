package service

import (
	"bookservice/interval/cache"
	"bookservice/interval/database/repository"
	"bookservice/pkg/models"
	"fmt"

	"github.com/google/uuid"
)

type BookService struct {
	bookRepo repository.IBookRepository
	cache    cache.IRedisClient
}

type IBookService interface {
	CreateBook(data *models.CreateBookRequestModel) error
	DeleteBookById(bookID uuid.UUID) error
	GetBookById(bookID uuid.UUID) (models.BookModel, error)
	GetAllBooks(queries map[string]string) ([]models.BookModel, error)
	UpdatePriceById(bookID uuid.UUID, newPrice string) error
}

func NewBookService() IBookService {
	return &BookService{
		bookRepo: repository.NewBookRepo(),
		cache:    cache.NewCacheClient(),
	}
}

func (bookService *BookService) CreateBook(data *models.CreateBookRequestModel) error {
	
	book, err := bookService.bookRepo.CreateBook(data)

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	
	bookService.cache.CreateBooks(book)
	
	return nil

}

func (bookService *BookService) DeleteBookById(bookID uuid.UUID) error {
	if err := bookService.bookRepo.DeleteBookById(bookID); err != nil {
		return fmt.Errorf(err.Error())
	}
	bookService.cache.DeleteBookById(bookID)
	return nil
}

func (bookService *BookService) GetBookById(bookID uuid.UUID) (models.BookModel, error) {
	cacheResult, err := bookService.cache.GetBookById(bookID)
	
	if err != nil || len(cacheResult) == 0 {
		books, err := bookService.bookRepo.GetBookById(bookID)
		if err != nil {
			return models.BookModel{}, fmt.Errorf(err.Error())
		}
		bookService.cache.CreateBooks([]models.BookModel{books})
		return books, nil
	}

	return cacheResult[0], nil
}

func (bookService *BookService) GetAllBooks(queries map[string]string) ([]models.BookModel, error) {

	cacheResult, err := bookService.cache.GetAllBooks(queries)
    
	if err != nil || len(cacheResult) == 0 {
		
		books, err := bookService.bookRepo.GetAllBooks(queries)
		if err != nil {
			return []models.BookModel{}, fmt.Errorf(err.Error())
		}
		bookService.cache.CreateBooks(books)
		return books, nil
	}

	return cacheResult, nil

}

func (bookService *BookService) UpdatePriceById(bookID uuid.UUID, newPrice string) error {
	if err := bookService.bookRepo.UpdatePriceById(bookID, newPrice); err != nil {
		return fmt.Errorf(err.Error())
	}
	bookService.cache.UpdatePriceById(bookID, newPrice)
	return nil
}
