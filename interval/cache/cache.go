package cache

import (
	"bookservice/pkg/models"
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisClient struct {
	rdb *redis.Client
}

type IRedisClient interface {
	CreateBooks(data []models.BookModel) error
	DeleteBookById(bookID uuid.UUID) error
	GetBookById(bookID uuid.UUID) ([]models.BookModel, error)
	GetAllBooks(queries map[string]string) ([]models.BookModel, error)
	UpdatePriceById(bookID uuid.UUID, newPrice string) error
}

func NewCacheClient() IRedisClient {
	return &RedisClient{
		rdb: redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDR"),
		}),
	}
}

func (redisCli *RedisClient) CreateBooks(books []models.BookModel) error {

	ctx := context.Background()
	wg := sync.WaitGroup{}
	for _, book := range books {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("book:%s", book.Id.String())
			redisCli.rdb.HSet(ctx, key, map[string]interface{}{
				"id":       book.Id.String(),
				"title":    book.Title,
				"author":   book.Author,
				"category": book.Category,
				"price":    book.Price,
			})
		}()
	}

	wg.Wait()

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
func (redisCli *RedisClient) GetBookById(bookID uuid.UUID) ([]models.BookModel, error) {

	ctx := context.Background()

	key := fmt.Sprintf("book:%s", bookID.String())

	result, err := redisCli.rdb.HGetAll(ctx, key).Result()

	if err != nil {
		return []models.BookModel{}, fmt.Errorf(err.Error())
	}

	if len(result) == 0 {
		return []models.BookModel{}, nil
	}

	var book models.BookModel

	book.Author = result["author"]
	book.Category = result["category"]
	book.Title = result["title"]
	book.Price = result["price"]

	created_at, _ := time.Parse(time.RFC3339, result["create_at"])

	book.Created_at = created_at

	id, _ := uuid.Parse(result["id"])

	book.Id = id

	return []models.BookModel{book}, nil
}
func (redisCli *RedisClient) GetAllBooks(queries map[string]string) ([]models.BookModel, error) {

	ctx := context.Background()

	query := ""

	for k, v := range queries {
		if k == "page" {
			continue
		}
		query += fmt.Sprintf("@%s:%s ", k, v)
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

	pageSize := start * 15

	result, err := redisCli.rdb.Do(ctx, "FT.SEARCH", "idx:books", query, "SORTBY", "created_at", "DESC", "LIMIT", strconv.Itoa(start), strconv.Itoa(pageSize)).Result()

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

			for i := 0; i < len(doc); i += 2 {

				switch doc[i].(string) {
				case "author":
					book.Author = doc[i+1].(string)
				case "category":
					book.Category = doc[i+1].(string)
				case "title":
					book.Title = doc[i+1].(string)
				case "price":
					book.Price = doc[i+1].(string)
				case "id":
					id, _ := uuid.Parse(doc[i+1].(string))
					book.Id = id
				case "create_at":
					created_at, _ := time.Parse(time.RFC3339, doc[i+1].(string))
					book.Created_at = created_at
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
