package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// 初始化 OpenTelemetry
func InitOTel(serviceName string) error {
	// 设置资源
	res, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %v", err)
	}

	// 初始化追踪
	tracerProvider, err := newTraceProvider(res)
	if err != nil {
		return fmt.Errorf("failed to create trace provider: %v", err)
	}

	// 初始化指标
	meterProvider, err := newMeterProvider(res)
	if err != nil {
		return fmt.Errorf("failed to create meter provider: %v", err)
	}

	// 设置全局实例
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)

	return traceProvider, nil
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	// Prometheus exporter
	prometheusExporter, err := prometheus.New(
		prometheus.WithoutScopeInfo(),
	)
	if err != nil {
		return nil, err
	}

	// 控制台导出器（用于调试）
	consoleExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	// 创建 MeterProvider
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(consoleExporter, metric.WithInterval(15*time.Second)),
		),
		metric.WithReader(
			prometheusExporter,
		),
		metric.WithResource(res),
	)

	return meterProvider, nil
}

// 关闭 OpenTelemetry 提供者
func ShutdownOTel(ctx context.Context) error {
	if tp, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
		if err := tp.Shutdown(ctx); err != nil {
			return err
		}
	}
	if mp, ok := otel.GetMeterProvider().(*metric.MeterProvider); ok {
		if err := mp.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

// 获取 Prometheus 指标处理器
func GetPrometheusHandler() (interface{}, error) {
	exporter, err := prometheus.New(
		prometheus.WithoutScopeInfo(),
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}
