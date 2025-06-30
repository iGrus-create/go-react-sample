package usecase

import (
	"os"
	"practice/model"
	"practice/repository"
	"practice/validator"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// ユーザーのユースケース
type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

// ユーザーのユースケースを作成
func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

// ユーザー登録
func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	// バリデーション
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Name: user.Name, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID: newUser.ID,
		Name: newUser.Name,
		Email: newUser.Email,
	}
	return resUser, nil
}

// ログイン
func (uu *userUsecase) Login(user model.User) (string, error) {
	// バリデーション
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}

	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
