package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"server_go/client"
	"server_go/proto"
)

// HelloServiceServer 实现 gRPC 服务端
type HelloServiceServer struct {
	proto.UnimplementedHelloServiceServer
	clientPool map[string]*client.HelloClient
}

func NewHelloServiceServer() *HelloServiceServer {
	return &HelloServiceServer{
		clientPool: make(map[string]*client.HelloClient),
	}
}

// SayHello 处理传入的 gRPC 请求（服务端功能）
func (s *HelloServiceServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	// 获取 tracer
	tracer := otel.Tracer("hello-service")

	// 开始一个 span
	ctx, span := tracer.Start(ctx, "SayHello")
	defer span.End()

	// 添加属性到 span
	span.SetAttributes(attribute.String("request.name", req.Name))

	log.Printf("Received gRPC request: %s", req.Name)

	// 模拟一些业务处理
	time.Sleep(10 * time.Millisecond)

	// 向其他微服务发起请求（客户端功能）
	responseFromOtherService, err := s.callOtherMicroservice(ctx, req)
	if err != nil {
		log.Printf("Warning: Failed to call other microservice: %v", err)
	}

	var combinedMessage string
	if responseFromOtherService != "" {
		combinedMessage = fmt.Sprintf("Hello %s! Your message: '%s'. Response from other service: %s",
			req.Name, req.Message, responseFromOtherService)
	} else {
		combinedMessage = fmt.Sprintf("Hello %s! Your message: '%s'", req.Name, req.Message)
	}

	response := &proto.HelloResponse{
		Message:   combinedMessage,
		Timestamp: time.Now().Unix(),
	}

	return response, nil
}

// 向其他微服务发起请求的内部方法（客户端功能）
func (s *HelloServiceServer) callOtherMicroservice(ctx context.Context, originalReq *proto.HelloRequest) (string, error) {
	// 检查是否已有到该服务的连接
	targetAddr := "other-microservice:50051" // 假设这是其他微服务的地址
	clientInstance, exists := s.clientPool[targetAddr]

	if !exists {
		// 创建新的客户端连接
		newClient, err := client.NewHelloClient(targetAddr)
		if err != nil {
			// 如果无法连接到其他微服务，返回错误但不中断当前请求
			log.Printf("Warning: Could not connect to other microservice %s: %v", targetAddr, err)
			return "", nil // Return empty string instead of error to not break the main request
		}
		clientInstance = newClient
		s.clientPool[targetAddr] = clientInstance
	}

	// 构造发往其他微服务的请求
	otherReq := &proto.HelloRequest{
		Name:    fmt.Sprintf("%s-via-proxy", originalReq.Name),
		Message: fmt.Sprintf("Forwarded from main service: %s", originalReq.Message),
	}

	// 调用其他微服务
	err := clientInstance.SayHelloWithContext(ctx, otherReq)
	if err != nil {
		return "", fmt.Errorf("failed to call other microservice: %v", err)
	}

	return fmt.Sprintf("Successfully contacted other service for %s", originalReq.Name), nil
}

// SayHelloStream 处理流式请求
func (s *HelloServiceServer) SayHelloStream(stream proto.HelloService_SayHelloStreamServer) error {
	tracer := otel.Tracer("hello-service")
	ctx, span := tracer.Start(stream.Context(), "SayHelloStream")
	defer span.End()

	for {
		req, err := stream.Recv()
		if err != nil {
			return status.Errorf(codes.Unknown, "unable to receive stream: %v", err)
		}

		log.Printf("Received stream request: %s", req.Name)

		// 向其他微服务转发请求
		_, err = s.callOtherMicroservice(ctx, req)
		if err != nil {
			log.Printf("Warning: Failed to forward to other microservice: %v", err)
		}

		response := &proto.HelloResponse{
			Message:   fmt.Sprintf("Hello %s! Your message: '%s'", req.Name, req.Message),
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(response); err != nil {
			return status.Errorf(codes.Unknown, "unable to send stream: %v", err)
		}
	}
}

// StartServer 启动 gRPC 服务器
func StartServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterHelloServiceServer(s, NewHelloServiceServer())

	log.Printf("Server listening on port %s", port)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
