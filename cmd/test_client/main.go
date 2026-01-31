package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"server_go/monitoring"
	"server_go/proto"
)

func main() {
	serviceName := "test-client"

	// 初始化 OpenTelemetry
	if err := monitoring.InitOTel(serviceName); err != nil {
		log.Printf("Warning: Failed to initialize OpenTelemetry: %v", err)
	}

	// 连接到我们的微服务
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	// 测试一元 RPC
	log.Println("Testing unary RPC...")
	req := &proto.HelloRequest{
		Name:    "TestClient",
		Message: "Hello from test client!",
	}

	resp, err := client.SayHello(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	log.Printf("Response received: %s", resp.Message)

	// 等待一段时间以确保指标被收集
	time.Sleep(2 * time.Second)

	log.Println("Test completed successfully!")
}
