package validator

import (
	"practice/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ITaskValidator interface {
	TaskValidate(task model.Task) error
}

type TaskValidator struct {}

func NewTaskValidator() ITaskValidator {
	return &TaskValidator{}
}

// タスクのバリデーション
func (tv *TaskValidator) TaskValidate(task model.Task) error {
	return validation.ValidateStruct(&task,
		validation.Field(
			&task.Title,
			validation.Required.Error("title is required"),
			validation.RuneLength(1, 10).Error("limited max 10 char"),
		),
	)
}