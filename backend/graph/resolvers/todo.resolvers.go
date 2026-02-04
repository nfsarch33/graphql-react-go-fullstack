package resolvers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nfsarch33/graphql-react-go-fullstack/backend/graph"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/graph/generated"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/graph/model"
	"github.com/nfsarch33/graphql-react-go-fullstack/backend/internal/models"
)

// Todos returns all todos with optional filtering
func (r *queryResolver) Todos(ctx context.Context, filter *model.TodoFilter) ([]*model.Todo, error) {
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
		todos[i] = &model.Todo{
			ID:          strconv.Itoa(t.ID),
			Title:       t.Title,
			Description: &t.Description,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	return todos, nil
}

// Todo returns a single todo by ID
func (r *queryResolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo id: %v", err)
	}

	var dbTodo models.Todo
	result := r.DB.First(&dbTodo, idInt)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, result.Error
	}

	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &dbTodo.Description,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// TodoCount returns total number of todos
func (r *queryResolver) TodoCount(ctx context.Context) (int, error) {
	var count int64
	result := r.DB.Model(&models.Todo{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// CreateTodo creates a new todo
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*model.Todo, error) {
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

	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &dbTodo.Description,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// UpdateTodo updates an existing todo
func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input model.UpdateTodoInput) (*model.Todo, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo id: %v", err)
	}

	var dbTodo models.Todo
	if err := r.DB.First(&dbTodo, idInt).Error; err != nil {
		return nil, fmt.Errorf("todo not found")
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

	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &dbTodo.Description,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// DeleteTodo deletes a todo
func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
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
func (r *mutationResolver) ToggleTodo(ctx context.Context, id string) (*model.Todo, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo id: %v", err)
	}

	var dbTodo models.Todo
	if err := r.DB.First(&dbTodo, idInt).Error; err != nil {
		return nil, fmt.Errorf("todo not found")
	}

	// Toggle completed status
	dbTodo.Completed = !dbTodo.Completed
	if err := r.DB.Save(&dbTodo).Error; err != nil {
		return nil, err
	}

	return &model.Todo{
		ID:          strconv.Itoa(dbTodo.ID),
		Title:       dbTodo.Title,
		Description: &dbTodo.Description,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}, nil
}

// Resolvers returns the resolver implementation
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
