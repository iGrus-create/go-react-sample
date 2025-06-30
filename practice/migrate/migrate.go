package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"practice/db"
	"practice/model"

	"github.com/aws/aws-lambda-go/lambda"
)

type MigrateResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func HandleMigrate(ctx context.Context) (MigrateResponse, error) {
	log.Printf("=== データベースマイグレーション開始 ===")
	
	// 環境変数を設定
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "production")
	}
	
	log.Printf("データベース接続を試行中...")
	dbConn := db.NewDB()
	if dbConn == nil {
		return MigrateResponse{
			Message: "データベース接続に失敗しました",
			Success: false,
		}, fmt.Errorf("database connection failed")
	}
	
	defer func() {
		if dbConn != nil {
			db.CloseDB(dbConn)
			log.Printf("データベース接続を閉じました")
		}
	}()
	
	log.Printf("テーブルマイグレーションを実行中...")
	log.Printf("User テーブルとTask テーブルを作成します...")
	
	// マイグレーション実行
	err := dbConn.AutoMigrate(&model.User{}, &model.Task{})
	if err != nil {
		log.Printf("マイグレーションエラー: %v", err)
		return MigrateResponse{
			Message: fmt.Sprintf("マイグレーションに失敗しました: %v", err),
			Success: false,
		}, err
	}
	
	log.Printf("マイグレーションが正常に完了しました")
	
	return MigrateResponse{
		Message: "データベースマイグレーションが正常に完了しました。UserテーブルとTaskテーブルが作成されました。",
		Success: true,
	}, nil
}

func main() {
	lambda.Start(HandleMigrate)
}
