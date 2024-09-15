package routes

import (
	"bookservice/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	bookController := controller.NewUserController()

	r.POST("/books", bookController.CreateBook)
	r.GET("/books", bookController.GetAllBooks)
	r.GET("/books/:id", bookController.GetBookById)
	r.DELETE("/books/:id", bookController.DeleteBookById)
	r.PATCH("/books/:id", bookController.UpdatePriceById)

	return r
}
