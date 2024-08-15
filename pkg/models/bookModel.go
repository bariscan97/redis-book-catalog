package models

import (
	"github.com/google/uuid"
)

type CreateBookRequestModel struct {
	Author   string `json:"author"   validate:"required"`
	Title    string `json:"title"    validate:"required"`
	Category string `json:"category" validate:"required"`
	Price    string `json:"price"    validate:"required"`
}


type BookModel struct {
	Id       uuid.UUID
	Author   string
	Title    string
	Category string
	Price    string
}
