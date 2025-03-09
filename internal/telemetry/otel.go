package telemetry

import (
	"context"
	"errors"
	"time"

	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
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

	loggerProvider, err := loggerProvider(ctx, conf, serviceName)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	hook := otellogrus.NewHook("qrquiz", otellogrus.WithLoggerProvider(loggerProvider))
	logrus.AddHook(hook)

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

func loggerProvider(ctx context.Context, conf *config.Config, serviceName string) (*log.LoggerProvider, error) {
	exporter, err := newLoggingExporter(ctx, conf)
	if err != nil {
		return nil, err
	}

	resource, err := defaultResource(serviceName)
	if err != nil {
		return nil, err
	}

	processor := log.NewBatchProcessor(exporter, log.WithExportInterval(time.Second*time.Duration(conf.Otlp.Interval)))
	loggingProvider := log.NewLoggerProvider(
		log.WithResource(resource),
		log.WithProcessor(processor),
	)
	return loggingProvider, nil
}

func newLoggingExporter(ctx context.Context, conf *config.Config) (log.Exporter, error) {
	var exporter log.Exporter
	var err error

	switch conf.Otlp.Protocol {

	case "grpc":
		var options []otlploggrpc.Option
		options = append(options, otlploggrpc.WithEndpoint(conf.Otlp.Endpoint))
		if conf.Otlp.Insecure {
			options = append(options, otlploggrpc.WithInsecure())
		}
		exporter, err = otlploggrpc.New(ctx, options...)

	default:
		exporter, err = stdoutlog.New()
	}

	return exporter, err
}
