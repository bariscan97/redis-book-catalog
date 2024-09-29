package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookRequestModel struct {
	Author   string `json:"author"   validate:"required"`
	Title    string `json:"title"    validate:"required"`
	Category string `json:"category" validate:"required"`
	Price    string `json:"price" binding:"required,numeric,gte=0"`
}

type BookModel struct {
	Id         uuid.UUID `json:"id"`
	Author     string    `json:"author"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	Price      string    `json:"price"`
	Created_at time.Time `json:"created_at"`
}
