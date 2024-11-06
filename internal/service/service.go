package service

import (
	"github.com/onfirebyte/todo-dumb/internal/model"
	"gorm.io/gorm"
)

type TodoService struct {
	db *gorm.DB
}

func NewTodoService(db *gorm.DB) (*TodoService, error) {
	return &TodoService{
		db: db,
	}, nil
}

func (t *TodoService) CreateTodo(todo model.Todo) error {
	return t.db.Create(&todo).Error
}

func (t *TodoService) GetTodosByUserId(userId uint) ([]model.Todo, error) {
	var todos []model.Todo
	err := t.db.Where("owner_id = ?", userId).Find(&todos).Error
	return todos, err
}

func (t *TodoService) UpdateTodoById(todo model.Todo) error {
	var old model.Todo
	err := t.db.First(&old, "id = ?", todo.ID).Error
	if err != nil {
		return err
	}

	old.Title = todo.Title
	old.Content = todo.Content
	old.Done = todo.Done

	return t.db.Save(&old).Error
}

func (t *TodoService) DeleteTodoById(todoId uint) error {
	return t.db.Delete(&model.Todo{}, todoId).Error
}

func (t *TodoService) IsOwner(todoId uint, userId uint) (bool, error) {
	var count int64
	err := t.db.Model(&model.Todo{}).Where("id = ? AND owner_id = ?", todoId, userId).Count(&count).Error
	return count > 0, err
}
