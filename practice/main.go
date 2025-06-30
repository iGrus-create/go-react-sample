package main

import (
	"fmt"
	"log"
	"os"
)

// startServer starts the HTTP server (for local development)
func startServer() {
	app := InitializeApp()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル用のデフォルト
	}
	
	log.Printf("=== サーバー設定 ===")
	log.Printf("ポート: %s", port)
	log.Printf("ポート %s でサーバーを起動中", port)
	
	// サーバー起動
	if err := app.Echo.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Printf("サーバーの起動に失敗しました: %v", err)
		os.Exit(1)
	}
}