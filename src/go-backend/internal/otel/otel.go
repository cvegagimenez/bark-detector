package otel

import (
	"context"
	"errors"
	"log"
	"os"
	"regexp"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	provider "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	meter metric.Meter

	stateMu       sync.RWMutex
	barkValueByID = map[string]float64{}
)

func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
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

	res, err := newResource()
	if err != nil {
		handleErr(err)
		return
	}

	meterProvider, err := newMeterProvider(ctx, res)
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)

	otel.SetMeterProvider(meterProvider)

	meter = otel.GetMeterProvider().Meter("bark-detector")

	err = setupMetrics()
	if err != nil {
		handleErr(err)
		return
	}

	return
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("bark-detector"),
		))
}

func newMeterProvider(ctx context.Context, res *resource.Resource) (*provider.MeterProvider, error) {
	metricExporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithEndpoint(envOrDefault("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "localhost:4318")),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`http\.server\.(request|response)\.size`)
	var dropMetricsView provider.View = func(i provider.Instrument) (provider.Stream, bool) {
		s := provider.Stream{
			Name:        i.Name,
			Description: i.Description,
			Unit:        i.Unit,
			Aggregation: provider.AggregationDrop{},
		}

		if re.MatchString(i.Name) {
			return s, true
		}

		return s, false
	}

	meterProvider := provider.NewMeterProvider(
		provider.WithResource(res),
		provider.WithReader(provider.NewPeriodicReader(metricExporter)),
		provider.WithView(
			dropMetricsView,
		),
	)

	return meterProvider, nil
}

func setupMetrics() error {
	barkGauge, err := meter.Float64ObservableGauge(
		"bark_power",
		metric.WithDescription("Latest bark power measurement received from the sensor"),
		metric.WithUnit("RMS"),
	)
	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			stateMu.RLock()
			defer stateMu.RUnlock()

			for sensorID, barkValue := range barkValueByID {
				observer.ObserveFloat64(
					barkGauge,
					barkValue,
					metric.WithAttributes(
						attribute.String("sensor.id", sensorID),
					),
				)
			}
			return nil
		},
		barkGauge,
	)
	if err != nil {
		return err
	}
	return nil
}

func RecordBarkPower(power float64, sensorID string) {
	stateMu.Lock()
	barkValueByID[sensorID] = power
	stateMu.Unlock()

	log.Printf("Recorded bark power: %f", power)
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
