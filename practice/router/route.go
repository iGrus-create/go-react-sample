package router

import (
	"net/http"
	"os"
	"practice/controller"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	e := echo.New()
	
	// 環境に応じたCORS設定
	var allowOrigins []string
	if os.Getenv("GO_ENV") == "production" || os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda環境では具体的なオリジンを指定（credentialsを使用するため*は使用不可）
		allowOrigins = []string{"https://d3ac4jqqf77oju.cloudfront.net"}
		// 環境変数で追加のオリジンを指定可能
		if feURL := os.Getenv("FE_URL"); feURL != "" {
			allowOrigins = append(allowOrigins, feURL)
		}
	} else {
		// ローカル開発環境では特定のオリジンを指定
		allowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
		if feURL := os.Getenv("FE_URL"); feURL != "" {
			allowOrigins = append(allowOrigins, feURL)
		}
	}
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowCredentials, echo.HeaderXCSRFToken},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))

	// ヘルスチェック用エンドポイント（CSRFチェック前に配置）
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"message": "Application is running",
		})
	})

	// CSRF設定（ヘルスチェックを除外）
	cookieDomain := os.Getenv("API_DOMAIN")
	if cookieDomain == "*" || cookieDomain == "" {
		cookieDomain = "" // 空文字にしてドメイン制限を無効化
	}
	
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath: "/",
		CookieDomain: cookieDomain,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
		CookieMaxAge: 60,
		Skipper: func(c echo.Context) bool {
			// ヘルスチェックエンドポイントはCSRFチェックをスキップ
			return c.Path() == "/health"
		},
	}))

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.Login)
	e.POST("/logout", uc.Logout)
	e.GET("/csrf", uc.CsrfToken)

	t := e.Group("/tasks")
	t.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	t.GET("", tc.GetAllTask)
	t.GET("/:task_id", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:task_id", tc.UpdateTask)
	t.DELETE("/:task_id", tc.DeleteTask)
	return e
}