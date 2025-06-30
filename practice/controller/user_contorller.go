package controller

import (
	"net/http"
	"os"
	"practice/model"
	"practice/usecase"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
    SignUp(c echo.Context) error
    Login(c echo.Context) error
    Logout(c echo.Context) error
	CsrfToken(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

// ユーザー登録
func (uc *userController) SignUp(c echo.Context) error {
	// リクエストボディをモデルにバインド
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ユースケース層を呼び出す
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// レスポンスを返す
	return c.JSON(http.StatusCreated, userRes)
}

// ログイン
func (uc *userController) Login(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// クッキーの設定
	cookie := new(http.Cookie)   // クッキー構造体を作成
	cookie.Name = "token"        // クッキーの名前
	cookie.Value = tokenString   // クッキーの値
	cookie.Expires = time.Now().Add(24 * time.Hour) // クッキーの有効期限
	cookie.Path = "/"            // クッキーのパス
	cookie.Domain = os.Getenv("API_DOMAIN")       // クッキーのドメイン
	cookie.Secure = true            // クッキーのセキュリティ
	cookie.HttpOnly = true            // クッキーのHttpOnly
	cookie.SameSite = http.SameSiteNoneMode // クッキーのSameSite
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

// ログアウト
// クッキーを削除
func (uc *userController) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

// CSRFトークン
func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}