package controller

import (
	"bookservice/internal/service"
	"bookservice/pkg/models"
	
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BookController struct {
	bookService service.IBookService
}

type IBookController interface {
	CreateBook(c *gin.Context)
	DeleteBookById(c *gin.Context)
	GetBookById(c *gin.Context)
	GetAllBooks(c *gin.Context)
	UpdatePriceById(c *gin.Context)
}

func NewUserController(bookservice service.IBookService) IBookController {
	return &BookController{
		bookService: bookservice,
	}
}

func (bookController *BookController) CreateBook(c *gin.Context) {

	validate := validator.New()

	var book models.CreateBookRequestModel

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validate.Struct(book); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err := strconv.Atoi(book.Price)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	if err := bookController.bookService.CreateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "succesful",
	})
}

func (bookController *BookController) DeleteBookById(c *gin.Context) {

	bookID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := bookController.bookService.DeleteBookById(bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "succesful",
	})
}

func (bookController *BookController) GetBookById(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	books, err := bookController.bookService.GetBookById(bookID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (bookController *BookController) GetAllBooks(c *gin.Context) {
	
	queries := make(map[string]string)
	
	queries["page"] = "0"
	
	ALL := c.Request.URL.Query()

	for i, j := range ALL {
		if len(j[0]) == 0 {
			continue
		}
		queries[i] = j[0]
		if i == "page" {
			if _, err := strconv.Atoi(j[0]); err != nil {
				queries[i] = "0"
			}
		}

	}

	books, err := bookController.bookService.GetAllBooks(queries)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (bookController *BookController) UpdatePriceById(c *gin.Context) {

	bookID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newPrice := c.Query("newprice")

	if _, err := strconv.Atoi(newPrice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := bookController.bookService.UpdatePriceById(bookID, newPrice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "succesful",
	})

}

