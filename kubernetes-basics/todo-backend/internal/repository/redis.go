package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"todo-backend/internal/models"
)

type RedisTodoRepository struct {
	client *redis.Client
}

func NewRedisTodoRepository() (*RedisTodoRepository, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		return nil, fmt.Errorf("REDIS_HOST is not set")
	}

	if port == "" {
		return nil, fmt.Errorf("REDIS_PORT is not set")
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	return &RedisTodoRepository{
		client: client,
	}, nil
}
func (r *RedisTodoRepository) InitializeFromJSON(filePath string) error {
	ctx := context.Background()

	// Check whether Redis already contains todos.
	exists, err := r.client.Exists(ctx, "todos").Result()
	if err != nil {
		return fmt.Errorf("failed to check Redis: %w", err)
	}

	// Database has already been initialized.
	if exists > 0 {
		return nil
	}

	// Read initial todos from JSON.
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open seed JSON: %w", err)
	}
	defer file.Close()

	var todos []models.Todo

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&todos); err != nil {
		return fmt.Errorf("failed to decode seed JSON: %w", err)
	}

	// Store every todo in Redis.
	for _, todo := range todos {
		data, err := json.Marshal(todo)
		if err != nil {
			return fmt.Errorf(
				"failed to encode todo %d: %w",
				todo.ID,
				err,
			)
		}

		key := fmt.Sprintf("todo:%d", todo.ID)

		if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
			return fmt.Errorf(
				"failed to save todo %d: %w",
				todo.ID,
				err,
			)
		}

		if err := r.client.SAdd(ctx, "todos", todo.ID).Err(); err != nil {
			return fmt.Errorf(
				"failed to index todo %d: %w",
				todo.ID,
				err,
			)
		}
	}

	// Make the next generated ID 12.
	if len(todos) > 0 {
		lastID := todos[len(todos)-1].ID

		if err := r.client.Set(
			ctx,
			"todo:next_id",
			lastID,
			0,
		).Err(); err != nil {
			return fmt.Errorf("failed to initialize next ID: %w", err)
		}
	}

	return nil
}
func (r *RedisTodoRepository) FetchAllTodos() ([]models.Todo, error) {
	ids, err := r.client.SMembers(context.Background(), "todos").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch todo IDs: %w", err)
	}

	todos := make([]models.Todo, 0, len(ids))

	for _, id := range ids {
		todoID, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("invalid todo ID %q: %w", id, err)
		}

		todo, err := r.FetchTodo(todoID)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}
func (r *RedisTodoRepository) FetchTodo(id int) (models.Todo, error) {
	key := fmt.Sprintf("todo:%d", id)

	data, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return models.Todo{}, fmt.Errorf("failed to fetch todo %d: %w", id, err)
	}

	var todo models.Todo

	if err := json.Unmarshal(data, &todo); err != nil {
		return models.Todo{}, fmt.Errorf(
			"failed to decode todo %d: %w",
			id,
			err,
		)
	}

	return todo, nil
}
func (r *RedisTodoRepository) CreateTodo(todo models.Todo) (models.Todo, error) {
	ctx := context.Background()

	// Generate a new ID
	id, err := r.client.Incr(ctx, "todo:next_id").Result()
	if err != nil {
		return models.Todo{}, fmt.Errorf("failed to generate todo ID: %w", err)
	}

	todo.ID = uint64(id)

	data, err := json.Marshal(todo)
	if err != nil {
		return models.Todo{}, fmt.Errorf("failed to encode todo: %w", err)
	}

	key := fmt.Sprintf("todo:%d", todo.ID)

	err = r.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		return models.Todo{}, fmt.Errorf("failed to save todo: %w", err)
	}

	err = r.client.SAdd(ctx, "todos", todo.ID).Err()
	if err != nil {
		return models.Todo{}, fmt.Errorf("failed to add todo ID to index: %w", err)
	}

	return todo, nil
}