package repository

import (
	"bookservice/internal/database"
	"bookservice/pkg/models"
	"context"
	"strings"

	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	BookRepository struct {
		pool *pgxpool.Pool
	}
	IBookRepository interface {
		CreateBook(data *models.CreateBookRequestModel) ([]models.BookModel, error)
		DeleteBookById(bookID uuid.UUID) error
		GetBookById(bookID uuid.UUID) (models.BookModel, error)
		GetAllBooks(queries map[string]string) ([]models.BookModel, error)
		UpdatePriceById(bookID uuid.UUID, newPrice string) error
	}
)

func NewBookRepo() IBookRepository {
	return &BookRepository{
		pool: database.Pool(),
	}
}

func (repo *BookRepository) CreateBook(data *models.CreateBookRequestModel) ([]models.BookModel, error) {
	ctx := context.Background()

	sql := `INSERT INTO books(title, author, category, price) VALUES($1, $2, $3 ,$4) RETURNING *`

	var book models.BookModel
	err := repo.pool.QueryRow(ctx, sql, data.Title, data.Author, data.Category, data.Price).Scan(&book.Id, &book.Title, &book.Author, &book.Category, &book.Price)
	if err != nil {
		return []models.BookModel{}, fmt.Errorf(err.Error())
	}

	return []models.BookModel{book}, nil
}

func (repo *BookRepository) DeleteBookById(bookID uuid.UUID) error {
	ctx := context.Background()

	sql := `
		DELETE FROM books 
		WHERE id = $1
	`
	_, err := repo.pool.Exec(ctx, sql, bookID)

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func (repo *BookRepository) GetBookById(bookID uuid.UUID) (models.BookModel, error) {
	ctx := context.Background()

	sql := `
		SELECT * FROM books
		WHERE id = $1
	`
	rows := repo.pool.QueryRow(ctx, sql, bookID)

	book := models.BookModel{}

	err := rows.Scan(&book.Id, &book.Author, &book.Title, &book.Category, &book.Price, &book.Created_at)

	if err != nil {
		return models.BookModel{}, err
	}

	return book, nil
}

func (repo *BookRepository) UpdatePriceById(bookID uuid.UUID, newPrice string) error {
	ctx := context.Background()

	sql := `
		UPDATE books
		SET
			price = $1
		WHERE
			id = $2
	`
	_, err := repo.pool.Exec(ctx, sql, newPrice, bookID)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func (repo *BookRepository) GetAllBooks(queries map[string]string) ([]models.BookModel, error) {
	ctx := context.Background()

	start_query := "SELECT * FROM books"

	addQueries := ""

	if len(queries) > 1 {
		addQueries += " WHERE "
	}

	page := queries["page"]

	parameters := make([]interface{}, 0)

	count := 1

	for key, val := range queries {
		if key == "page" {
			continue
		}
		addQueries += " " + key + " = " + fmt.Sprintf("$%d", count)
		parameters = append(parameters, val)

		if count < len(queries) {
			addQueries += " AND"
		}
		count += 1
	}

	if strings.HasSuffix(addQueries, " AND") {
		addQueries = addQueries[:len(addQueries)-4]
	}

	parameters = append(parameters, page)

	order := " ORDER BY create_at DESC "

	pagination := fmt.Sprintf(" LIMIT 15 OFFSET $%d * 15", count)

	sql := start_query + addQueries + order + pagination
	rows, err := repo.pool.Query(ctx, sql, parameters...)

	if err != nil {
		return []models.BookModel{}, err
	}

	var Books []models.BookModel

	for rows.Next() {
		var book models.BookModel

		if DbError := rows.Scan(&book.Id, &book.Author, &book.Title, &book.Category, &book.Price, &book.Created_at); DbError != nil {
			return []models.BookModel{}, fmt.Errorf(DbError.Error())
		}

		Books = append(Books, book)
	}

	return Books, nil
}
