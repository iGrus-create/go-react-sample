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

type InspectResponse struct {
	Message    string                 `json:"message"`
	Success    bool                   `json:"success"`
	Tables     []string               `json:"tables"`
	UserCount  int64                  `json:"user_count"`
	TaskCount  int64                  `json:"task_count"`
	Users      []model.UserResponse   `json:"users"`
	Tasks      []model.Task           `json:"tasks"`
}

func HandleInspect(ctx context.Context) (InspectResponse, error) {
	log.Printf("=== データベース確認開始 ===")
	
	// 環境変数を設定
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "production")
	}
	
	log.Printf("データベース接続を試行中...")
	dbConn := db.NewDB()
	if dbConn == nil {
		return InspectResponse{
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
	
	response := InspectResponse{
		Success: true,
		Tables:  []string{},
		Users:   []model.UserResponse{},
		Tasks:   []model.Task{},
	}
	
	// テーブル一覧を取得
	var tableNames []string
	err := dbConn.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tableNames).Error
	if err != nil {
		log.Printf("テーブル一覧取得エラー: %v", err)
	} else {
		response.Tables = tableNames
		log.Printf("テーブル一覧: %v", tableNames)
	}
	
	// ユーザー数を取得
	var userCount int64
	err = dbConn.Model(&model.User{}).Count(&userCount).Error
	if err != nil {
		log.Printf("ユーザー数取得エラー: %v", err)
	} else {
		response.UserCount = userCount
		log.Printf("ユーザー数: %d", userCount)
	}
	
	// タスク数を取得
	var taskCount int64
	err = dbConn.Model(&model.Task{}).Count(&taskCount).Error
	if err != nil {
		log.Printf("タスク数取得エラー: %v", err)
	} else {
		response.TaskCount = taskCount
		log.Printf("タスク数: %d", taskCount)
	}
	
	// ユーザー一覧を取得（最大10件）
	var users []model.User
	err = dbConn.Select("id, name, email, created_at, updated_at").Limit(10).Find(&users).Error
	if err != nil {
		log.Printf("ユーザー一覧取得エラー: %v", err)
	} else {
		for _, user := range users {
			response.Users = append(response.Users, model.UserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			})
		}
		log.Printf("ユーザー一覧を取得しました（%d件）", len(response.Users))
	}
	
	// タスク一覧を取得（最大10件）
	err = dbConn.Limit(10).Find(&response.Tasks).Error
	if err != nil {
		log.Printf("タスク一覧取得エラー: %v", err)
	} else {
		log.Printf("タスク一覧を取得しました（%d件）", len(response.Tasks))
	}
	
	response.Message = fmt.Sprintf("データベース確認完了。テーブル数: %d, ユーザー数: %d, タスク数: %d", 
		len(response.Tables), response.UserCount, response.TaskCount)
	
	log.Printf("=== データベース確認完了 ===")
	return response, nil
}

func main() {
	lambda.Start(HandleInspect)
} 