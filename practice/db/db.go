package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
    if os.Getenv("GO_ENV") == "dev" {
        err := godotenv.Load()
        if err != nil {
            log.Printf("警告: .envファイルが見つかりません: %v", err)
        }
    }
    
    var url string
    if os.Getenv("GO_ENV") == "production" {
        // 本番環境（RDS PostgreSQL）- 接続タイムアウトを追加
        url = os.Getenv("DATABASE_URL")
        if url == "" {
            // DATABASE_URLが設定されていない場合は個別の環境変数から構築
            url = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require&connect_timeout=10",
                os.Getenv("POSTGRES_USER"),
                os.Getenv("POSTGRES_PW"),
                os.Getenv("POSTGRES_HOST"),
                os.Getenv("POSTGRES_PORT"),
                os.Getenv("POSTGRES_DB"))
        }
    } else {
        // 開発環境（ローカルPostgreSQL）
        url = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=5",
            os.Getenv("POSTGRES_USER"),
            os.Getenv("POSTGRES_PW"),
            os.Getenv("POSTGRES_HOST"),
            os.Getenv("POSTGRES_PORT"),
            os.Getenv("POSTGRES_DB"))
    }
    
    log.Printf("データベースに接続中...")
    
    // GORM設定（ログレベルを調整）
    config := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent), // ログを抑制
    }
    
    db, err := gorm.Open(postgres.Open(url), config)
    if err != nil {
        log.Printf("データベース接続に失敗しました: %v", err)
        return nil
    }
    
    // 接続プールの設定
    sqlDB, err := db.DB()
    if err == nil {
        sqlDB.SetMaxIdleConns(2)
        sqlDB.SetMaxOpenConns(10)
        sqlDB.SetConnMaxLifetime(time.Hour)
    }
    
    log.Printf("データベース接続に成功しました")
    return db
}

func CloseDB(db *gorm.DB) {
    if db == nil {
        return
    }
    sqlDB, _ := db.DB()
    if err := sqlDB.Close(); err != nil {
        log.Printf("データベース接続のクローズに失敗しました: %v", err)
    }
}
