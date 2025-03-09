package telemetry

import (
	"context"
	"errors"
	"time"

	"github.com/sekthor/qrquiz/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func SetUpTelemetry(ctx context.Context, conf *config.Config, serviceName string) (shutdown func(context.Context) error, err error) {

	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	tracerProvider, err := tracerProvider(ctx, conf, serviceName)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	return
}

func defaultResource(serviceName string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			// TODO: inject the service version
			// semconv.ServiceVersion(serviceVersion)
		),
	)
}

func tracerProvider(ctx context.Context, conf *config.Config, serviceName string) (*trace.TracerProvider, error) {
	traceExporter, err := newTraceExporter(ctx, conf)
	if err != nil {
		return nil, err
	}

	resource, err := defaultResource(serviceName)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithResource(resource),
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(
				time.Duration(conf.Otlp.Interval)*time.Second)),
	)
	return traceProvider, nil
}

func newTraceExporter(ctx context.Context, conf *config.Config) (trace.SpanExporter, error) {
	var exporter trace.SpanExporter
	var err error

	switch conf.Otlp.Protocol {

	case "grpc":
		var options []otlptracegrpc.Option
		options = append(options, otlptracegrpc.WithEndpoint(conf.Otlp.Endpoint))
		if conf.Otlp.Insecure {
			options = append(options, otlptracegrpc.WithInsecure())
		}
		exporter, err = otlptracegrpc.New(ctx, options...)

	default:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	return exporter, err
}
