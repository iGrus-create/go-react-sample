package repository

import (
	"errors"
	"fmt"
	"practice/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITaskRepository interface {
	GetAllTask(tasks *[]model.Task, userId uint) error
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

// taskの構造体を作成
type taskRepository struct {
	db *gorm.DB
}

// taskのコンストラクタ
func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

// タスクを全件取得
func (tr *taskRepository) GetAllTask(tasks *[]model.Task, userId uint) error {
	if tr.db == nil {
		return errors.New("database not connected")
	}
	if err := tr.db.Joins("User").Where("user_id = ?", userId).Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

// タスクを1件取得
func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	if tr.db == nil {
		return errors.New("database not connected")
	}
	if err := tr.db.Joins("User").Where("user_id = ?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

// タスクを作成
func (tr *taskRepository) CreateTask(task *model.Task) error {
	if tr.db == nil {
		return errors.New("database not connected")
	}
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

// タスクの更新
func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
	if tr.db == nil {
		return errors.New("database not connected")
	}
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object dose not exist")
	}
	return nil
}

// タスクの削除
func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	if tr.db == nil {
		return errors.New("database not connected")
	}
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object dose not exist")
	}
	return nil
}
