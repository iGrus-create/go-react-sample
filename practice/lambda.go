package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

var echoLambda *echoadapter.EchoLambda

// Handler AWS Lambdaのメインハンドラー関数
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("=== Lambdaハンドラー実行 ===")
	log.Printf("リクエストパス: %s", req.Path)
	log.Printf("リクエストメソッド: %s", req.HTTPMethod)
	
	// Echoアダプターを使用してリクエストを処理
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Lambda環境かどうかの判定（複数の条件で確認）
	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" || 
		       os.Getenv("AWS_EXECUTION_ENV") != "" ||
		       os.Getenv("_LAMBDA_SERVER_PORT") != ""
	
	if isLambda {
		// Running in AWS Lambda
		log.Printf("=== Lambda初期化 ===")
		log.Printf("Lambda変数名: %s", os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))
		log.Printf("Lambda実行環境: %s", os.Getenv("AWS_EXECUTION_ENV"))
		log.Printf("Lambdaポート: %s", os.Getenv("_LAMBDA_SERVER_PORT"))
		
		// 環境変数の読み込み
		if os.Getenv("GO_ENV") == "" {
			os.Setenv("GO_ENV", "production")
		}
		
		// アプリケーション初期化
		app := InitializeApp()
		
		// Lambda用のアダプターを作成
		echoLambda = echoadapter.New(app.Echo)
		
		log.Printf("Lambdaハンドラーの初期化が正常に完了しました")
		
		lambda.Start(Handler)
	} else {
		// Running locally
		log.Printf("ローカル環境で実行中")
		startServer()
	}
} 