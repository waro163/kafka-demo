package common

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TracingConfig struct {
	CollectorHost string  //tracing_host
	SamplingRate  float64 //tracing_sample_rate
}

// var tp *sdktrace.TracerProvider
// tp, err := InitTracer("your-app-name", xxx)
// tp.Shutdown(context.Background())

func InitTracer(serverName string, tracingCfg TracingConfig) (*sdktrace.TracerProvider, error) {
	exporter, err := initTraceExporter(tracingCfg.CollectorHost)
	if err != nil {
		return nil, err
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(tracingCfg.SamplingRate))),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(serverName))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

func initTraceExporter(host string) (sdktrace.SpanExporter, error) {
	var traceExporter sdktrace.SpanExporter
	var err error
	if strings.Contains(host, "4317") {
		traceExporter, err = initOtelGrpcTracerExporter(host)
	} else if strings.Contains(host, "4318") {
		traceExporter, err = initOtelHttpTracerExporter(host)
	} else if strings.Contains(host, "14268") {
		traceExporter, err = initJaegerExporter(host)
	} else {
		traceExporter, err = initStdoutExporter(os.Stdout)
	}
	return traceExporter, err
}

func initOtelGrpcTracerExporter(host string) (sdktrace.SpanExporter, error) {
	// host should contain :4317 port
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create otel grpc trace exporter: %w", err)
	}
	return traceExporter, nil
}

func initOtelHttpTracerExporter(host string) (sdktrace.SpanExporter, error) {
	// host should contain :4318 port
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(host), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to create otel http trace exporter: %w", err)
	}
	return traceExporter, nil
}

func initJaegerExporter(host string) (sdktrace.SpanExporter, error) {
	// jaeger collector should be like http://localhost:14268/api/traces
	traceExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(host)))
	return traceExporter, err
}

func initStdoutExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint(), stdouttrace.WithWriter(w))
	return traceExporter, err
}
