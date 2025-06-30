package main

import (
	"log"
	"os"
	"practice/controller"
	"practice/db"
	"practice/repository"
	"practice/router"
	"practice/usecase"
	"practice/validator"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// App holds the application dependencies
type App struct {
	Echo *echo.Echo
	DB   *gorm.DB
}

// InitializeApp sets up the application with all dependencies
func InitializeApp() *App {
	log.Printf("=== アプリケーション初期化開始 ===")
	log.Printf("環境: %s", os.Getenv("GO_ENV"))
	
	// 一時的にデータベース接続をスキップ（テスト用）
	var database *gorm.DB = nil
	if os.Getenv("SKIP_DB") != "true" {
		// データベース接続
		log.Printf("データベース接続を試行中...")
		database = db.NewDB()
		if database == nil {
			log.Printf("警告: データベース接続に失敗しましたが、データベースなしで続行します")
		} else {
			log.Printf("データベース接続に成功しました")
		}
	} else {
		log.Printf("データベース接続をスキップしました（テストモード）")
	}

	log.Printf("バリデーターを初期化中...")
	userValidator := validator.NewUserValidator()
	taskValidator := validator.NewTaskValidator()

	log.Printf("リポジトリを初期化中...")
	userRepository := repository.NewUserRepository(database)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	userController := controller.NewUserController(userUsecase)

	taskRepository := repository.NewTaskRepository(database)
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	taskController := controller.NewTaskController(taskUsecase)

	log.Printf("ルーターを設定中...")
	e := router.NewRouter(userController, taskController)
	
	log.Printf("=== アプリケーション設定 ===")
	log.Printf("環境: %s", os.Getenv("GO_ENV"))
	log.Printf("データベース接続: %v", database != nil)
	log.Printf("アプリケーション初期化が正常に完了しました")
	
	return &App{
		Echo: e,
		DB:   database,
	}
} 