package main

import (
	"context"
	"log"
	"time"

	"server_go/client"
	"server_go/monitoring"
)

func main() {
	serviceName := "hello-client"

	// 初始化 OpenTelemetry
	if err := monitoring.InitOTel(serviceName); err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	// 延迟关闭 OpenTelemetry
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := monitoring.ShutdownOTel(ctx); err != nil {
			log.Printf("Error shutting down OpenTelemetry: %v", err)
		}
	}()

	// 创建客户端
	client, err := client.NewHelloClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// 测试一元 RPC
	log.Println("Testing unary RPC...")
	if err := client.SayHello("World", "This is a test message"); err != nil {
		log.Printf("Error calling SayHello: %v", err)
	}

	// 短暂延迟
	time.Sleep(1 * time.Second)

	// 测试流式 RPC
	log.Println("Testing streaming RPC...")
	names := []string{"Alice", "Bob", "Charlie"}
	ctx := context.Background()
	if err := client.SayHelloStream(ctx, names); err != nil {
		log.Printf("Error calling SayHelloStream: %v", err)
	}

	log.Println("Client finished")
}
