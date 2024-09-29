package main

import (
	"bookservice/internal/cache"
	"bookservice/internal/controller"
	"bookservice/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), 
	})

	
	redisCli := cache.NewCacheClient(client)

	bookService := service.NewBookService(redisCli)

	bookController := controller.NewUserController(bookService)

	app := gin.Default()

	bookRoutes := app.Group("/books")
	{
		bookRoutes.POST("/", bookController.CreateBook)       
		bookRoutes.GET("/", bookController.GetAllBooks)       
		bookRoutes.GET("/:id", bookController.GetBookById)    
		bookRoutes.PUT("/:id", bookController.GetBookById) 
		bookRoutes.DELETE("/:id", bookController.DeleteBookById)  
	}

	log.Fatal(app.Run(os.Getenv("PORT")))
}
