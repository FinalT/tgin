package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"testing"
)

func TestOpenTelemetry(t *testing.T) {
	ctx := context.Background()
	tracer := otel.GetTracerProvider().
		Tracer("github.com/FinalT/tgin/extra/opentelemetry")

	// 如果 ctx 已经和一个 span 绑定了，那么新的 span 就是老的 span 的日子
	ctx, span := tracer.Start(ctx, "opentelemetry-demo",
		trace.WithAttributes(attribute.String("version", "1")))
	defer span.End()

	// 重置名字
	span.SetName("otel-demo")
	span.SetAttributes(attribute.Int("status", 200))
	span.AddEvent("马老师，发生什么事了")

}

func initZipkin(t *testing.T) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		t.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func initJeager(t *testing.T) {
	url := "http://localhost:14268/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		t.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)

	otel.SetTracerProvider(tp)
}
