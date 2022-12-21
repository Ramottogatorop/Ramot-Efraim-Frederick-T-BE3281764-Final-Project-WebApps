package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&tasks).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return tasks, nil // TODO: replace this
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	Task := r.db.WithContext(ctx).Create(&task).Error
	if Task != nil {
		return 0, Task
	}
	return task.ID, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	task := entity.Task{}
	err := r.db.WithContext(ctx).First(&task, id)
	return task, err.Error // TODO: replace this
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	task := []entity.Task{}
	err := r.db.WithContext(ctx).Where("category_id = ?", catId).Find(&task)
	if err.Error != nil {
		return []entity.Task{}, err.Error
	}
	return task, err.Error
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id =?", task.ID).Updates(&task).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	var tasks []entity.Task
	err := r.db.WithContext(ctx).Delete(&tasks, id).Error
	if err != nil {
		return err
	}
	return nil
}
