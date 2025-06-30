package usecase

import (
	"practice/model"
	"practice/repository"
	"practice/validator"
)

type ITaskUsecase interface {
	GetAllTask(userId uint) ([]model.TaskResponse, error)
	GetTaskById(userId uint, taskId uint) (model.TaskResponse, error)
	CreateTask(ask model.Task) (model.TaskResponse, error)
	UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error)
	DeleteTask(userId uint, taskId uint) error
}

type taskUsecase struct {
	tr repository.ITaskRepository
	tv validator.ITaskValidator
}

func NewTaskUsecase(tr repository.ITaskRepository, tv validator.ITaskValidator) ITaskUsecase {
	return &taskUsecase{tr, tv}
}

// タスクを全件取得
func (tu *taskUsecase) GetAllTask(userId uint) ([]model.TaskResponse, error) {
	tasks := []model.Task{}
	if err := tu.tr.GetAllTask(&tasks, userId); err != nil {
		return nil, err
	}
	resTasks := []model.TaskResponse{}
	for _, task := range tasks {
		t := model.TaskResponse{
			ID: task.ID,
			Title: task.Title,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		}
		resTasks = append(resTasks, t)
	}
	return resTasks, nil
}

// タスクを一件取得
func (tu *taskUsecase) GetTaskById(userId uint, taskId uint) (model.TaskResponse, error) {
	task := model.Task{}
	if err := tu.tr.GetTaskById(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}
	resTask := model.TaskResponse{
		ID: task.ID,
		Title: task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}


// タスクを作成
func (tu *taskUsecase) CreateTask(task model.Task) (model.TaskResponse, error) {
	// バリデーション
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}

	// タスクの作成
	if err := tu.tr.CreateTask(&task); err != nil {
		return model.TaskResponse{}, err
	}
	resTask := model.TaskResponse{
		ID: task.ID,
		Title: task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

// タスクの更新
func (tu *taskUsecase) UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error) {
	// バリデーション
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}

	// タスクの更新
	if err := tu.tr.UpdateTask(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}

	// レスポンス
	resTask := model.TaskResponse{
		ID: task.ID,
		Title: task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

// タスクを削除
func (ut *taskUsecase) DeleteTask(userId uint, taskId uint) error {
	if err := ut.tr.DeleteTask(userId, taskId); err != nil {
		return err
	}
	return nil
}


