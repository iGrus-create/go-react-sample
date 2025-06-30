package validator

import (
	"practice/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
)

type IUserValidator interface {
	UserValidate(user model.User) error
}

type UserValidator struct {}

func NewUserValidator() IUserValidator {
	return &UserValidator{}
}

// ユーザーのバリデーション
func (uv *UserValidator) UserValidate(user model.User) error {
	return validation.ValidateStruct(&user,
		validation.Field(
			&user.Email,
			validation.Required.Error("メールアドレスは必須項目です。"),
			validation.RuneLength(1, 30).Error("メールアドレスは30文字以内で入力してください。"),
			is.Email.Error("メールアドレスの形式が正しくありません。"),
		),
		validation.Field(
			&user.Password,
			validation.Required.Error("パスワードは必須項目です。"),
			validation.RuneLength(8, 30).Error("パスワードは8文字以上30文字以内で入力してください。"),
		),
	)
}