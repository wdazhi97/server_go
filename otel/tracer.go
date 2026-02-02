package otel

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracer 初始化 OpenTelemetry tracer
func InitTracer(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	// 获取 OTLP endpoint，默认使用 otel-collector
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "otel-collector:4317"
	}

	// 创建 gRPC 连接
	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	// 创建 OTLP exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// 创建 resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", os.Getenv("ENVIRONMENT")),
		),
	)
	if err != nil {
		return nil, err
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 设置全局 propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// GinMiddleware 返回一个 Gin 中间件用于追踪 HTTP 请求
func GinMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		// 从请求头提取 trace context
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// 创建 span
		spanName := c.Request.Method + " " + c.FullPath()
		if c.FullPath() == "" {
			spanName = c.Request.Method + " " + c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethod(c.Request.Method),
				semconv.HTTPURL(c.Request.URL.String()),
				semconv.HTTPRoute(c.FullPath()),
				semconv.NetHostName(c.Request.Host),
				attribute.String("http.client_ip", c.ClientIP()),
			),
		)
		defer span.End()

		// 获取 trace ID 和 span ID
		spanCtx := span.SpanContext()
		traceID := spanCtx.TraceID().String()
		spanID := spanCtx.SpanID().String()

		// 将 trace ID 添加到响应头（方便调试）
		c.Header("X-Trace-ID", traceID)
		c.Header("X-Span-ID", spanID)

		// 将 trace ID 存储到 Gin context（方便在 handler 中使用）
		c.Set("trace_id", traceID)
		c.Set("span_id", spanID)

		// 记录开始时间
		start := time.Now()

		// 将 context 注入到请求中
		c.Request = c.Request.WithContext(ctx)

		// 处理请求
		c.Next()

		// 记录响应信息
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		span.SetAttributes(
			semconv.HTTPStatusCode(statusCode),
			attribute.Int64("http.response_time_ms", duration.Milliseconds()),
		)

		// 如果有错误，记录错误
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("error", c.Errors.String()))
		}

		// 打印请求日志（包含 Trace ID，方便在服务器端追踪）
		logLevel := "INFO"
		if statusCode >= 400 && statusCode < 500 {
			logLevel = "WARN"
		} else if statusCode >= 500 {
			logLevel = "ERROR"
		}

		log.Printf("[%s] [trace_id=%s] %s %s | status=%d | duration=%dms | client_ip=%s",
			logLevel,
			traceID,
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			duration.Milliseconds(),
			c.ClientIP(),
		)
	}
}

// GetTracer 获取指定服务的 tracer
func GetTracer(serviceName string) trace.Tracer {
	return otel.Tracer(serviceName)
}

// StartSpan 开始一个新的 span
func StartSpan(ctx context.Context, tracer trace.Tracer, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}

// GetTraceID 从 context 中获取 trace ID
func GetTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

// GetSpanID 从 context 中获取 span ID
func GetSpanID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}

// GetTraceIDFromGin 从 Gin context 中获取 trace ID
func GetTraceIDFromGin(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	return ""
}

// InjectContext 将 trace context 注入到 gRPC metadata
func InjectContext(ctx context.Context) context.Context {
	return ctx
}
