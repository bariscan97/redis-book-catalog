package routes

import (
    "github.com/gin-gonic/gin"
    "bookservice/interval/controller"
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