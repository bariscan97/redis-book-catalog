package cache

import (
	"bookservice/models"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisClient struct {
	rdb *redis.Client
}

type IRedisClient interface {
	CreateBooks(data *models.CreateBookRequestModel) error
	DeleteBookById(bookID uuid.UUID) error
	GetBookById(bookID uuid.UUID) (*models.BookModel, error)
	GetAllBooks(queries map[string]string) ([]models.BookModel, error)
	UpdatePriceById(bookID uuid.UUID, newPrice string) error
}

func NewCacheClient(redisCli *redis.Client) IRedisClient {
	return &RedisClient{
		rdb: redisCli,
	}
}

func (redisCli *RedisClient) CreateBooks(book *models.CreateBookRequestModel) error {

	ctx := context.Background()

	ID := uuid.New()

	key := fmt.Sprintf("book:%s", ID)

	err := redisCli.rdb.HSet(ctx, key, map[string]interface{}{
		"id":         ID.String(),
		"title":      book.Title,
		"author":     book.Author,
		"category":   book.Category,
		"price":      book.Price,
		"created_at": time.Now().Unix(),
	})

	if err.Err() != redis.Nil {
		return err.Err()
	}

	return nil
}
func (redisCli *RedisClient) DeleteBookById(bookID uuid.UUID) error {

	ctx := context.Background()

	key := fmt.Sprintf("book:%s", bookID.String())

	_, err := redisCli.rdb.Del(ctx, key).Result()

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
func (redisCli *RedisClient) GetBookById(bookID uuid.UUID) (*models.BookModel, error) {

	ctx := context.Background()

	key := fmt.Sprintf("book:%s", bookID.String())

	result, err := redisCli.rdb.HGetAll(ctx, key).Result()

	if err != nil {
		return &models.BookModel{}, fmt.Errorf(err.Error())
	}

	var book models.BookModel
	ID, err := uuid.Parse(result["id"])
	if err != nil {
		return &models.BookModel{}, fmt.Errorf(err.Error()) 
	}
	book.Id = ID
	book.Author = result["author"]
	book.Category = result["category"]
	book.Title = result["title"]
	book.Price = result["price"]

	tmx, err := strconv.ParseInt(result["created_at"], 10, 64)

	if err != nil {
		panic(err)
	}

	tm := time.Unix(tmx, 0)
	book.Created_at = tm

	return &book, nil
}
func (redisCli *RedisClient) GetAllBooks(queries map[string]string) ([]models.BookModel, error) {

	ctx := context.Background()

	query := ""

	for key, value := range queries {

		if key == "page" {
			continue
		}

		switch key {
		case "gt":
			num, err := strconv.Atoi(value)
			if err != nil {
				return []models.BookModel{}, fmt.Errorf(err.Error())
			}
			query += fmt.Sprintf("@%s:[%d +inf] ", key, num)

		case "lt":
			num, err := strconv.Atoi(value)
			if err != nil {
				return []models.BookModel{}, fmt.Errorf(err.Error())
			}
			query += fmt.Sprintf("@%s:[0 %d] ", key, num)

		default:
			query += fmt.Sprintf("@%s:%s ", key, value)
		}
	}

	if query == "" {
		query = "*"
	}

	var start int

	num, err := strconv.Atoi(queries["page"])

	if err != nil {
		start = 0
	} else {
		start = num
	}

	pageSize := 15

	offset := pageSize * start

	result, err := redisCli.rdb.Do(ctx, "FT.SEARCH", "idx:books", query, "SORTBY", "created_at", "DESC", "LIMIT", offset, pageSize).Result()

	if err != nil {
		return []models.BookModel{}, nil
	}

	resultArray, ok := result.([]interface{})

	if !ok {
		return []models.BookModel{}, nil
	}

	CacheResults := []models.BookModel{}

	for i := 1; i < len(resultArray); i++ {

		doc, ok := resultArray[i].([]interface{})

		var book models.BookModel

		if ok {

			for j := 0; j < len(doc); j += 2 {

				switch doc[j].(string) {
				case "id" :
					ID, err := uuid.Parse(doc[j + 1].(string))
					if err != nil {
						break
					}
					book.Id = ID
				case "author":
					book.Author = doc[j+1].(string)
				case "category":
					book.Category = doc[j+1].(string)
				case "title":
					book.Title = doc[j+1].(string)
				case "price":
					book.Price = doc[j+1].(string)
				case "created_at":
					tmx, err := strconv.ParseInt(doc[j+1].(string), 10, 64)
					if err != nil {
						break
					}
					tm := time.Unix(tmx, 0)
					book.Created_at = tm
				}

			}

			CacheResults = append(CacheResults, book)
		}

	}

	return CacheResults, nil
}

func (redisCli *RedisClient) UpdatePriceById(bookID uuid.UUID, newPrice string) error {

	ctx := context.Background()

	key := fmt.Sprintf("book:%s", bookID.String())

	field := "price"

	_, err := redisCli.rdb.HSet(ctx, key, field, newPrice).Result()

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}
