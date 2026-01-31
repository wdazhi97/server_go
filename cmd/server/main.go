package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"server_go/monitoring"
	"server_go/server"
)

func main() {
	serviceName := "hello-service"

	// 初始化 OpenTelemetry
	if err := monitoring.InitOTel(serviceName); err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	// 注册信号处理以便优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// 启动 gRPC 服务器
	log.Println("Starting gRPC server (serves requests and calls other microservices)...")
	if err := server.StartServer("50051"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// 关闭 OpenTelemetry
	if err := monitoring.ShutdownOTel(ctx); err != nil {
		log.Printf("Error shutting down OpenTelemetry: %v", err)
	}
}
