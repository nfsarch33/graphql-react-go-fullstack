package graph

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/nfsarch33/graphql-react-go-fullstack/backend/graph/generated"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/graph/model"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/internal/models"
	"gorm.io/gorm"
)

// Todos returns all todos with optional filtering
func (r *Resolver) Todos(ctx context.Context, filter *model.TodoFilter) ([]*model.Todo, error) {
	var dbTodos []models.Todo
	query := r.DB

	// Apply filters if provided
	if filter != nil {
		if filter.Completed != nil {
			query = query.Where("completed = ?", *filter.Completed)
		}
		if filter.Search != nil && *filter.Search != "" {
			search := fmt.Sprintf("%%%s%%", *filter.Search)
			query = query.Where("title LIKE ? OR description LIKE ?", search, search)
		}
	}

	result := query.Order("created_at DESC").Find(&dbTodos)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert database models to GraphQL models
	todos := make([]*model.Todo, len(dbTodos))
	for i, t := range dbTodos {
		desc := t.Description // avoid taking address of range variable field directly
		todos[i] = &model.Todo{
			ID:          strconv.Itoa(t.ID),
			Title:       t.Title,
			Description: &desc,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	return todos, nil
}

// Todo returns a single todo by ID
func (r *Resolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo id: %v", err)
	}

	var dbTodo models.Todo
	result := r.DB.First(&dbTodo, idInt)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, result.Error
	}

	desc := dbTodo.Description
	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &desc,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// TodoCount returns total number of todos
func (r *Resolver) TodoCount(ctx context.Context) (int, error) {
	var count int64
	result := r.DB.Model(&models.Todo{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// CreateTodo creates a new todo
func (r *Resolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*model.Todo, error) {
	dbTodo := models.Todo{
		Title:       input.Title,
		Description: "",
		Completed:   false,
	}

	if input.Description != nil {
		dbTodo.Description = *input.Description
	}

	result := r.DB.Create(&dbTodo)
	if result.Error != nil {
		return nil, result.Error
	}

	desc := dbTodo.Description
	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &desc,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// UpdateTodo updates an existing todo
func (r *Resolver) UpdateTodo(ctx context.Context, id string, input model.UpdateTodoInput) (*model.Todo, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo id: %v", err)
	}

	var dbTodo models.Todo
	if err := r.DB.First(&dbTodo, idInt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, err
	}

	// Update only provided fields
	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Completed != nil {
		updates["completed"] = *input.Completed
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	if err := r.DB.Model(&dbTodo).Updates(updates).Error; err != nil {
		return nil, err
	}

	desc := dbTodo.Description
	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &desc,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// DeleteTodo deletes a todo
func (r *Resolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("invalid todo id: %v", err)
	}

	result := r.DB.Delete(&models.Todo{}, idInt)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

// ToggleTodo toggles the completion status of a todo
func (r *Resolver) ToggleTodo(ctx context.Context, id string) (*model.Todo, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo id: %v", err)
	}

	var dbTodo models.Todo
	if err := r.DB.First(&dbTodo, idInt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, err
	}

	// Toggle completed status
	dbTodo.Completed = !dbTodo.Completed
	if err := r.DB.Save(&dbTodo).Error; err != nil {
		return nil, err
	}

	desc := dbTodo.Description
	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &desc,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// Query returns the QueryResolver implementation
func (r *Resolver) Query() generated.QueryResolver { return r }

// Mutation returns the MutationResolver implementation
func (r *Resolver) Mutation() generated.MutationResolver { return r }
