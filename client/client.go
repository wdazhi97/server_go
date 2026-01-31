package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"server_go/proto"
)

type HelloClient struct {
	conn    *grpc.ClientConn
	client  proto.HelloServiceClient
	address string
}

func NewHelloClient(address string) (*HelloClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	client := proto.NewHelloServiceClient(conn)

	log.Printf("Connected to remote service at %s", address)

	return &HelloClient{
		conn:    conn,
		client:  client,
		address: address,
	}, nil
}

func (c *HelloClient) Close() error {
	return c.conn.Close()
}

// SayHelloWithContext 允许传递上下文（用于传播追踪信息）
func (c *HelloClient) SayHelloWithContext(ctx context.Context, req *proto.HelloRequest) error {
	// 使用 OpenTelemetry 追踪
	ctx, span := otel.Tracer("hello-client").Start(ctx, "SayHello")
	defer span.End()

	resp, err := c.client.SayHello(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to call SayHello: %v", err)
	}

	log.Printf("Response received from %s: %s (timestamp: %d)", c.address, resp.Message, resp.Timestamp)
	return nil
}

// SayHello 简单的 SayHello 方法
func (c *HelloClient) SayHello(name, message string) error {
	ctx := context.Background()
	req := &proto.HelloRequest{
		Name:    name,
		Message: message,
	}

	return c.SayHelloWithContext(ctx, req)
}

// SayHelloStream 发送流式请求
func (c *HelloClient) SayHelloStream(ctx context.Context, names []string) error {
	stream, err := c.client.SayHelloStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to call SayHelloStream: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Printf("Stream response received from %s: %s (timestamp: %d)", c.address, in.Message, in.Timestamp)
		}
	}()

	for _, name := range names {
		err := stream.Send(&proto.HelloRequest{
			Name:    name,
			Message: fmt.Sprintf("Hello from %s!", name),
		})
		if err != nil {
			return fmt.Errorf("failed to send stream: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	stream.CloseSend()
	<-waitc

	return nil
}

// Address 返回客户端连接的地址
func (c *HelloClient) Address() string {
	return c.address
}
